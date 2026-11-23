#!/bin/bash
# Quản lý tiến trình Go Backend & Tự động hồi sinh (Auto-restart)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR" || exit 1

ACTION="${1:-restart}"
export APP_ENV=staging

case "$ACTION" in
    run-loop)
        # Tiến trình giám sát ngầm: chạy ./server và tự động restart nếu bị tắt đột ngột
        echo "[$(date '+%Y-%m-%d %H:%M:%S')] [SUPERVISOR] Khởi động supervisor..." >> app.log
        while true; do
            if ! pgrep -f "./server" > /dev/null 2>&1; then
                echo "[$(date '+%Y-%m-%d %H:%M:%S')] [SUPERVISOR] Khởi chạy ./server..." >> app.log
                ./server >> app.log 2>&1
                EXIT_CODE=$?
                echo "[$(date '+%Y-%m-%d %H:%M:%S')] [SUPERVISOR] CẢNH BÁO: ./server đã dừng với exit code ${EXIT_CODE}. Tự động khởi động lại sau 2 giây..." >> app.log
            fi
            sleep 2
        done
        ;;

    start)
        if pgrep -f "./server" > /dev/null 2>&1; then
            echo "===> [Server] ./server đã đang chạy với PID: $(pgrep -f "./server" | head -n 1)"
            exit 0
        fi
        echo "===> [Server] Đang khởi chạy supervisor trong background..."
        chmod +x ./server service.sh 2>/dev/null || true
        nohup bash "$0" run-loop > /dev/null 2>&1 &
        SUPERVISOR_PID=$!
        sleep 2
        SERVER_PID=$(pgrep -f "./server" | head -n 1)
        echo "===> [Server] Supervisor PID: ${SUPERVISOR_PID}"
        echo "===> [Server] Go Server PID: ${SERVER_PID}"
        echo "===> [Server] Khởi động thành công! Cơ chế tự phục hồi đã sẵn sàng."
        ;;

    stop)
        echo "===> [Server] Đang dừng supervisor và server..."
        pkill -9 -f "run-loop" || true
        pkill -9 -f "./server" || true
        sleep 2
        killall -9 server 2>/dev/null || true
        echo "===> [Server] Đã dừng toàn bộ tiến trình."
        ;;

    restart)
        echo "===> [Server] Đang dừng supervisor và server cũ..."
        pkill -9 -f "run-loop" || true
        pkill -9 -f "./server" || true
        sleep 2
        killall -9 server 2>/dev/null || true
        sleep 1

        echo "===> [Server] Đang khởi chạy supervisor mới trong background..."
        chmod +x ./server service.sh 2>/dev/null || true
        nohup bash "$0" run-loop > /dev/null 2>&1 &
        SUPERVISOR_PID=$!
        sleep 2
        SERVER_PID=$(pgrep -f "./server" | head -n 1)
        echo "===> [Server] Supervisor PID: ${SUPERVISOR_PID}"
        echo "===> [Server] Go Server PID: ${SERVER_PID}"
        echo "===> [Server] Server đã restart thành công và tự động hồi sinh khi bị crash!"
        echo "===> [Server] Xem log bằng lệnh: tail -f app.log"
        ;;

    status)
        echo "=== Trạng thái Backend ==="
        if pgrep -f "./server" > /dev/null 2>&1; then
            echo "[OK] ./server đang chạy với PID: $(pgrep -f "./server" | tr '\n' ' ')"
        else
            echo "[STOPPED] ./server không chạy!"
        fi

        if pgrep -f "service.sh run-loop" > /dev/null 2>&1; then
            echo "[OK] Supervisor đang chạy với PID: $(pgrep -f "service.sh run-loop" | tr '\n' ' ')"
        else
            echo "[WARN] Supervisor không chạy!"
        fi
        ;;

    migrate)
        echo "===> [Server] Đang kiểm tra và áp dụng Database Migrations..."
        DB_USER=$(grep -E '^DB_USER=' .env | cut -d '=' -f2- | tr -d ' "\r\n')
        DB_PASS=$(grep -E '^DB_PASS=' .env | cut -d '=' -f2- | tr -d ' "\r\n')
        DB_HOST=$(grep -E '^DB_HOST=' .env | cut -d '=' -f2- | tr -d ' "\r\n')
        DB_PORT=$(grep -E '^DB_PORT=' .env | cut -d '=' -f2- | tr -d ' "\r\n')
        DB_NAME=$(grep -E '^DB_NAME=' .env | cut -d '=' -f2- | tr -d ' "\r\n')

        [ -z "$DB_HOST" ] && DB_HOST="127.0.0.1"
        [ -z "$DB_PORT" ] && DB_PORT="3306"

        mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "CREATE TABLE IF NOT EXISTS schema_migrations (version VARCHAR(128) NOT NULL PRIMARY KEY, applied_at DATETIME DEFAULT CURRENT_TIMESTAMP);" 2>/dev/null

        for sql_file in $(ls -1 db/migrations/*.up.sql 2>/dev/null | sort); do
            version=$(basename "$sql_file" | cut -d '_' -f1)
            applied=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -N -s -e "SELECT count(*) FROM schema_migrations WHERE version='$version';" 2>/dev/null)
            if [ "$applied" = "0" ]; then
                echo "--> Đang chạy migration: $(basename "$sql_file")..."
                mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" < "$sql_file"
                if [ $? -eq 0 ]; then
                    mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "INSERT INTO schema_migrations (version) VALUES ('$version');" 2>/dev/null
                    echo "--> Áp dụng thành công: $(basename "$sql_file")"
                else
                    echo "--> [LỖI] Migration thất bại: $(basename "$sql_file")"
                    exit 1
                fi
            fi
        done
        echo "===> [Server] Database migrations hoàn tất!"
        echo "===> [Server] Đang chạy seeder master data..."
        chmod +x ./server 2>/dev/null || true
        ./server seed-master || true
        ;;

    seed-master|seed-master-data)
        TABLE="$2"
        if [ -n "$TABLE" ]; then
            echo "===> [Server] Đang chạy seeder master data cho bảng: ${TABLE}..."
            chmod +x ./server 2>/dev/null || true
            ./server seed-master "$TABLE"
        else
            echo "===> [Server] Đang chạy seeder toàn bộ master data..."
            chmod +x ./server 2>/dev/null || true
            ./server seed-master
        fi
        ;;

    check)
        # Dành cho cron job: nếu server chưa chạy thì tự gọi restart
        if ! pgrep -f "./server" > /dev/null 2>&1; then
            echo "[$(date '+%Y-%m-%d %H:%M:%S')] [WATCHDOG] Server chưa chạy. Đang tự động kích hoạt..." >> app.log
            bash "$0" restart >> app.log 2>&1
        fi
        ;;

    *)
        echo "Cách sử dụng: $0 {start|stop|restart|status|migrate|seed-master [table]|check}"
        exit 1
        ;;
esac

