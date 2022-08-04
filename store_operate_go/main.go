package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	USERNAME = "root"
	PASSWORD = "root"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "store"
)

type Clothes struct {
	clothing_id    string `json:"clothing_id " form:"clothing_id "`
	clothing_size  string `json:"clothing_size" form:"clothing_size"`
	clothing_price int    `json:"clothing_price" form:"clothing_price"`
	clothing_type  string `json:"clothing_type  " form:"clothing_type  "` // 0 正常状态， 1删除
}

type depository struct {
	depository_id       string `json:"depository_id" form:"depository_id"`
	depository_capacity string `json:"depository_capacity" form:"depository_capacity"`
}

type supplier struct {
	supplier_id   string `json:"supplier___id" form:"supplier_id"`
	supplier_name string `json:"supplier___name" form:"supplier_name"`
}

func main() {
	DB := Init()
	Insert1(DB)
	Query1(DB)
	Query2(DB)
	Query3(DB)
	Query4(DB)
	Query5(DB)
	Update(DB)
	Del(DB)
}

func Init() (DB *sql.DB) {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	DB, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println("connection to mysql failed:", err)
		return
	} else {
		fmt.Println("connection to mysql successed:", err)
		DB.SetMaxOpenConns(100)
		DB.SetMaxIdleConns(10)
		return
	}
}

func Insert1(DB *sql.DB) {
	q := []string{"insert into depository values ('001', 10), ('002', 20);\n",
		"insert into clothes values ('Blue01', 'S', 99, 'dress'), ('Red01', 'S', 100, 'dress');\n",
		"insert into supplier values ('001', 'tencent'), ('002', 'wangyi');\n",
		"insert into supply_condition values ('Blue01', '001', 'S+'), ('Red01', '002', 'D');\n"}
	for i := 0; i < len(q); i++ {
		_, err := DB.Exec(q[i])
		if err != nil {
			println("插入失败")
			//log.Fatal(err)
		}
		println("插入成功")
	}
}

func Query1(DB *sql.DB) {
	println("(1)查询服装尺码为'S'且销售价格在100以下的服装信息。")
	q := `
	select *
	from clothes
	where clothing_size = 'S' && clothes.clothing_price < 100;
`
	rows, err := DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		Temp := new(Clothes)
		if err := rows.Scan(&Temp.clothing_size, &Temp.clothing_id, &Temp.clothing_price, &Temp.clothing_type); err != nil {
			log.Fatal(err)
		}
		log.Printf("size %s id %s price %d type %s\n", Temp.clothing_size, Temp.clothing_id, Temp.clothing_price, Temp.clothing_type)
	}
}

func Query2(DB *sql.DB) {
	println("(2)查询仓库容量最大的仓库信息。")
	q := `
	select *
	from depository
	where depository_capacity = (
    	select max(depository_capacity)
    	from depository
    );
`
	rows, err := DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		Temp := new(depository)
		if err := rows.Scan(&Temp.depository_id, &Temp.depository_capacity); err != nil {
			log.Fatal(err)
		}
		log.Printf("id %s capacity %s\n", Temp.depository_id, Temp.depository_capacity)
	}
}

func Query3(DB *sql.DB) {
	println("（3）查询A类服装的库存总量。")
	q := `
	select depository_capacity
	from depository
	where depository_id = any(
    	select supplier_id
    	from supply_condition
    	where clothing_id = any(
        	select clothing_id
        	from clothes
        	where clothing_type = 'dress'
        )
	);
`
	rows, err := DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		Temp := new(depository)
		if err := rows.Scan(&Temp.depository_capacity); err != nil {
			log.Fatal(err)
		}
		log.Printf("capacity %s\n", Temp.depository_capacity)
	}
}

func Query4(DB *sql.DB) {
	println("(4) 查询服装编码以‘B’开始开头的服装。")
	q := `
	select *
	from clothes
	where clothing_id like 'B%';
`

	rows, err := DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		Temp := new(Clothes)
		if err := rows.Scan(&Temp.clothing_size, &Temp.clothing_id, &Temp.clothing_price, &Temp.clothing_type); err != nil {
			log.Fatal(err)
		}
		log.Printf("size %s id %s price %d type %s\n", Temp.clothing_size, Temp.clothing_id, Temp.clothing_price, Temp.clothing_type)
	}
}

func Query5(DB *sql.DB) {
	println("（5）查询服装质量等级有不合格的供应商信息。")
	q := `
	select *
	from supplier
	where supplier_id = (
    	select supplier_id
    	from supply_condition
    	where clothing_level <= 'D'
    );
`

	rows, err := DB.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		Temp := new(supplier)
		if err := rows.Scan(&Temp.supplier_id, &Temp.supplier_name); err != nil {
			log.Fatal(err)
		}
		log.Printf("size %s id %s price %d type %s\n", Temp.supplier_id, Temp.supplier_name)
	}
}

func Update(DB *sql.DB) {
	q := `
	update clothes
	set clothing_price = clothing_price*1.1;
`
	_, err := DB.Exec(q)
	if err != nil {
		println("更新失败")
		log.Fatal(err)
	}
	println("更新成功")
}

func Del(DB *sql.DB) {
	q := `
delete from supply_condition
where clothing_level <= 'D';
`
	_, err := DB.Exec(q)
	if err != nil {
		println("删除失败")
		log.Fatal(err)
	}
	println("删除成功")
}
