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

    check)
        # Dành cho cron job: nếu server chưa chạy thì tự gọi restart
        if ! pgrep -f "./server" > /dev/null 2>&1; then
            echo "[$(date '+%Y-%m-%d %H:%M:%S')] [WATCHDOG] Server chưa chạy. Đang tự động kích hoạt..." >> app.log
            bash "$0" restart >> app.log 2>&1
        fi
        ;;

    *)
        echo "Cách sử dụng: $0 {start|stop|restart|status|check}"
        exit 1
        ;;
esac

