package seeder

import (
	"fmt"
	"log"
	"strings"

	"platform-go/internal/infrastructure/repository"
)

// MasterSeederFunc defines the signature for a master table seeder function.
type MasterSeederFunc func(db repository.DBInterface) error

// MasterTable defines a master data table registration.
type MasterTable struct {
	Name        string
	Description string
	Seeder      MasterSeederFunc
}

// MasterDataTables is the central registry of all master data tables.
// Sau này nếu có bất kỳ bảng nào là master data:
// CHỈ CẦN THÊM TÊN BẢNG VÀ HÀM SEED TƯƠNG ỨNG VÀO DANH SÁCH DƯỚI ĐÂY:
var MasterDataTables = []MasterTable{
	{
		Name:        "travel_themes",
		Description: "Master data for travel themes (30 default themes)",
		Seeder: func(db repository.DBInterface) error {
			themes := CreateTravelThemes(db)
			log.Printf("[MASTER DATA] 'travel_themes': %d records ready.", len(themes))
			return nil
		},
	},
	// Ví dụ sau này có thêm bảng master data mới, chỉ cần thêm 1 entry như sau:
	// {
	// 	Name:        "locations",
	// 	Description: "Master data for locations",
	// 	Seeder: func(db repository.DBInterface) error {
	// 		CreateLocations(db)
	// 		return nil
	// 	},
	// },
}

// GetMasterTableNames returns the list of all registered master table names.
func GetMasterTableNames() []string {
	names := make([]string, 0, len(MasterDataTables))
	for _, t := range MasterDataTables {
		names = append(names, t.Name)
	}
	return names
}

// FindMasterTable finds a master table by name (case-insensitive).
func FindMasterTable(tableName string) (*MasterTable, bool) {
	for i := range MasterDataTables {
		if strings.EqualFold(MasterDataTables[i].Name, tableName) {
			return &MasterDataTables[i], true
		}
	}
	return nil, false
}

// RunMasterSeeds executes master data seeders.
// - If targetTables is empty: executes ALL registered master tables in order.
// - If targetTables contains specific table names: executes only those master tables.
func RunMasterSeeds(db repository.DBInterface, targetTables ...string) error {
	var tablesToRun []MasterTable

	if len(targetTables) == 0 {
		tablesToRun = MasterDataTables
		log.Printf("[MASTER SEED] Starting seed for all %d registered master table(s): %v", len(tablesToRun), GetMasterTableNames())
	} else {
		for _, rawName := range targetTables {
			name := strings.TrimSpace(rawName)
			if name == "" {
				continue
			}
			table, found := FindMasterTable(name)
			if !found {
				return fmt.Errorf("unknown master table '%s'. Available master tables: %v", name, GetMasterTableNames())
			}
			tablesToRun = append(tablesToRun, *table)
		}
		log.Printf("[MASTER SEED] Starting seed for %d targeted master table(s)...", len(tablesToRun))
	}

	for i, table := range tablesToRun {
		log.Printf("[MASTER SEED] (%d/%d) Seeding master table '%s'...", i+1, len(tablesToRun), table.Name)
		if err := table.Seeder(db); err != nil {
			return fmt.Errorf("failed to seed master table '%s': %w", table.Name, err)
		}
	}

	log.Printf("[MASTER SEED] Successfully completed master data seeding.")
	return nil
}
