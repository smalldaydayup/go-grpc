package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.comgo-grpc-http-rest-microservice-tutorial/pkg/protocol/grpc"
	v1 "github.comgo-grpc-http-rest-microservice-tutorial/pkg/service/v1"

	_ "github.com/go-sql-driver/mysql"
)

// Config是Server的配置
type Config struct {
	// gRPC服務器啟動引數部分
	// GRPCPort是gRPC服務器監聽的TCP端口
	GRPCPort string

	// 資料庫資料儲存引數部分
	// DatestoreDBHost是資料庫的地址
	DatastoreDBHost string
	// DatastoreDBUser是用於連接資料庫的用戶名
	DatastoreDBUser string
	// DatastoreDBPassword是用於連接資料庫的密碼
	DatastoreDBPassword string
	// DatastoreDBSchema是資料庫的名稱
	DatastoreDBSchema string
}

// RunServer運行gRPC服務器和HTTP網關
func RunServer() error {
	ctx := context.Background()

	// 獲取配置
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// 添加MySQL驅動程式特定引數來解析 date/time
	// 為另一個資料庫刪除它
	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword, cfg.DatastoreDBHost, cfg.DatastoreDBSchema, param)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewToDoServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
