package scheduler

import (
	"CronDashboardIBS/database"
	"CronDashboardIBS/functions"
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/shopspring/decimal"
	"github.com/vjeantet/jodaTime"
	"log"
	_ "math/big"
	"net/http"
	"os"
	"strconv"
	_ "strings"
	"time"
)

type SelectDashboardReal struct {
	RealisasiToday          float64 `json:"realisasi_today"`
	RealisasiTodayCount     int     `json:"realisasi_today_count"`
	RealisasiMonth          float64 `json:"realisasi_month"`
	RealisasiMonthCount     int     `json:"realisasi_month_count"`
	BukaTabunganToday       float64 `json:"buka_tabungan_today"`
	BukaTabunganTodayCount  int     `json:"buka_tabungan_today_count"`
	BukaTabunganMonth       float64 `json:"buka_tabungan_month"`
	BukaTabunganMonthCount  int     `json:"buka_tabungan_month_count"`
	BukaDepositoToday       float64 `json:"buka_deposito_today"`
	BukaDepositoTodayCount  int     `json:"buka_deposito_today_count"`
	BukaDepositoMonth       float64 `json:"buka_deposito_month"`
	BukaDepositoMonthCount  int     `json:"buka_deposito_month_count"`
	TransAngsuranPokok      float64 `json:"trans_angsuran_pokok"`
	TransAngsuranBunga      float64 `json:"trans_angsuran_bunga"`
	TransAngsuranDenda      float64 `json:"trans_angsuran_denda"`
	TransAngsuranCount      int     `json:"trans_angsuran_count"`
	TransSetorTabungan      float64 `json:"trans_setor_tabungan"`
	TransSetorTabunganCount int     `json:"trans_setor_tabungan_count"`
	TransTarikTabungan      float64 `json:"trans_tarik_tabungan"`
	TransTarikTabunganCount int     `json:"trans_tarik_tabungan_count"`
	TransSetorDeposito      float64 `json:"trans_setor_deposito"`
	TransSetorDepositoCount int     `json:"trans_setor_deposito_count"`
	TransTarikDeposito      float64 `json:"trans_tarik_deposito"`
	TransTarikDepositoCount int     `json:"trans_tarik_deposito_count"`
	KodeKantor              string  `json:"kode_kantor"`
}

type SelectDashboardChart struct {
	Period     string  `json:"period"`
	Amount     float64 `json:"amount"`
	Amount2    float64 `json:"amount2"`
	Amount3    float64 `json:"amount3"`
	Amount4    float64 `json:"amount4"`
	Amount5    float64 `json:"amount5"`
	Count      float64 `json:"count"`
	KodeKantor string  `json:"kode_kantor"`
}

type SelectDashboardUser struct {
	UserId      int    `json:"user_id"`
	UserName    string `json:"user_name"`
	NamaLengkap string `json:"nama_lengkap"`
	Jabatan     string `json:"jabatan"`
	Flag        int    `json:"flag"`
	Menu        string `json:"menu"`
}

type SelectDashboardReport struct {
	Period     string  `json:"period"`
	GorD       string  `json:"g_or_d"`
	Level      int     `json:"level"`
	NamaPerk   string  `json:"nama_perk"`
	KodePerk   string  `json:"kode_perk"`
	KodeInduk  string  `json:"kode_induk"`
	TypePerk   string  `json:"type_perk"`
	Amount     float64 `json:"amount"`
	Amount2    float64 `json:"amount2"`
	Amount3    float64 `json:"amount3"`
	Amount4    float64 `json:"amount4"`
	Amount5    float64 `json:"amount5"`
	KodeKantor string  `json:"kode_kantor"`
}

type SelectDashboardReportKredit struct {
	Period        string  `json:"period"`
	Kode          string  `json:"kode"`
	MyKolekNumber string  `json:"keterangan"`
	BakiDebet     float64 `json:"saldo"`
	JmlRek        int     `json:"jml_rek"`
	KodeKantor    string  `json:"kode_kantor"`
}

type SelectDashboardReportTKS struct {
	NoID      string  `json:"no_id"`
	SandiPos  string  `json:"sandi_pos"`
	GOrD      string  `json:"g_or_d"`
	Deskripsi string  `json:"deskripsi"`
	Jumlah    float64 `json:"jumlah"`
	Tanggal   string  `json:"tanggal"`
}

var db = database.ConnectDB()

func GetDataDashboardReal() {

	functions.Logger().Info("Starting Scheduler Get Data Dashboard Real ")

	cTgl := jodaTime.Format("YYYYMMdd", time.Now())
	cMonth := jodaTime.Format("MM", time.Now())
	cYear := jodaTime.Format("YYYY", time.Now())
	RealisasiToday := 0.00
	RealisasiTodayCount := 0
	RealisasiMonth := 0.00
	RealisasiMonthCount := 0

	BukaTabunganToday := 0.00
	BukaTabunganTodayCount := 0
	BukaTabunganMonth := 0.00
	BukaTabunganMonthCount := 0

	BukaDepositoToday := 0.00
	BukaDepositoTodayCount := 0
	BukaDepositoMonth := 0.00
	BukaDepositoMonthCount := 0

	TransAngsuranPokok := 0.00
	TransAngsuranBunga := 0.00
	TransAngsuranDenda := 0.00
	TransAngsuranCount := 0

	TransSetorTabungan := 0.00
	TransSetorTabunganCount := 0
	TransTarikTabungan := 0.00
	TransTarikTabunganCount := 0

	TransSetorDeposito := 0.00
	TransSetorDepositoCount := 0
	TransTarikDeposito := 0.00
	TransTarikDepositoCount := 0

	sqlStatement := "SELECT kode_kantor from app_kode_kantor order by kode_kantor "
	rowsX, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}
	for rowsX.Next() {
		messageList := SelectDashboardReal{}
		err = rowsX.Scan(&messageList.KodeKantor)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		KodeKantor := messageList.KodeKantor

		//REALISASI HARI INI
		sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiToday,COUNT(m.no_rekening) RealisasiTodayCount " +
			"FROM " +
			"kretrans t," +
			"nasabah n," +
			"kredit m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND m.tgl_realisasi = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=1 and pokok>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.RealisasiToday, &messageList.RealisasiTodayCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			RealisasiToday = messageList.RealisasiToday
			RealisasiTodayCount = messageList.RealisasiTodayCount

		}

		//END OF REALISASI HARI INI

		//REALISASI BULAN INI
		sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
			"FROM " +
			"kretrans t," +
			"nasabah n," +
			"kredit m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND MONTH(m.tgl_realisasi) = '" + cMonth + "' " +
			"AND YEAR(m.tgl_realisasi) = '" + cYear + "' " +
			"and floor(my_kode_trans/100)=1  and pokok>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.RealisasiMonth, &messageList.RealisasiMonthCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			RealisasiMonth = messageList.RealisasiMonth
			RealisasiMonthCount = messageList.RealisasiMonthCount
		}

		//END OF REALISASI BULAN INI

		//BUKA TABUNGAN HARI INI
		sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
			"FROM " +
			"tabtrans t," +
			"nasabah n," +
			"tabung m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND m.kode_produk not in ('ABP','AKP','KPE') " +
			"AND (m.tgl_register) = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=1 and pokok>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.BukaTabunganToday, &messageList.BukaTabunganTodayCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			BukaTabunganToday = messageList.BukaTabunganToday
			BukaTabunganTodayCount = messageList.BukaTabunganTodayCount
		}

		//END OF BUKA TABUNGAN HARI INI

		//BUKA TABUNGAN BULAN INI
		sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
			"FROM " +
			"tabtrans t," +
			"nasabah n," +
			"tabung m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND m.kode_produk not in ('ABP','AKP','KPE') " +
			"AND MONTH(m.tgl_register) = '" + cMonth + "' " +
			"AND YEAR(m.tgl_register) = '" + cYear + "' " +
			"and floor(my_kode_trans/100)=1 and pokok>0 and m.kode_kantor='" + KodeKantor + "'  "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.BukaTabunganMonth, &messageList.BukaTabunganMonthCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			BukaTabunganMonth = messageList.BukaTabunganMonth
			BukaTabunganMonthCount = messageList.BukaTabunganMonthCount
		}

		//END OF BUKA TABUNGAN BULAN INI

		//BUKA DEPOSITO HARI INI
		sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
			"FROM " +
			"deptrans t," +
			"nasabah n," +
			"deposito m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND (m.tgl_registrasi) = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=1 and pokok_trans>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.BukaDepositoToday, &messageList.BukaDepositoTodayCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			BukaDepositoToday = messageList.BukaDepositoToday
			BukaDepositoTodayCount = messageList.BukaDepositoTodayCount
		}

		//END OF BUKA DEPOSITO HARI INI

		//BUKA DEPOSITO BULAN INI
		sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
			"FROM " +
			"deptrans t," +
			"nasabah n," +
			"deposito m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND MONTH(m.tgl_registrasi) = '" + cMonth + "' " +
			"AND YEAR(m.tgl_registrasi) = '" + cYear + "' " +
			"and floor(my_kode_trans/100)=1 and pokok_trans>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.BukaDepositoMonth, &messageList.BukaDepositoMonthCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			BukaDepositoMonth = messageList.BukaDepositoMonth
			BukaDepositoMonthCount = messageList.BukaDepositoMonthCount
		}

		//END OF BUKA DEPOSITO BULAN INI

		//TRANS ANGSURAN
		sqlStatement = "SELECT coalesce(SUM(pokok),0) principal, coalesce(SUM(bunga),0) interest, coalesce(SUM(denda),0) fines, COUNT(t.no_rekening) count " +
			"FROM " +
			"kretrans t," +
			"nasabah n," +
			"kredit m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND (t.tgl_trans) = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=3 and (pokok<>0 or bunga<>0 or denda<>0) and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.TransAngsuranPokok, &messageList.TransAngsuranBunga,
				&messageList.TransAngsuranDenda, &messageList.TransAngsuranCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			TransAngsuranPokok = messageList.TransAngsuranPokok
			TransAngsuranBunga = messageList.TransAngsuranBunga
			TransAngsuranDenda = messageList.TransAngsuranDenda
			TransAngsuranCount = messageList.TransAngsuranCount
		}

		//END OF TRANS ANGSURAN

		//TRANS SETORAN TABUNGAN
		sqlStatement = "SELECT coalesce(SUM(pokok),0) principal, COUNT(t.no_rekening) count " +
			"FROM " +
			"tabtrans t," +
			"nasabah n," +
			"tabung m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND (t.tgl_trans) = '" + cTgl + "' " +
			"AND m.kode_produk not in ('ABP','AKP','KPE') " +
			"and floor(my_kode_trans/100)=1 and pokok>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.TransSetorTabungan, &messageList.TransSetorTabunganCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			TransSetorTabungan = messageList.TransSetorTabungan
			TransSetorTabunganCount = messageList.TransSetorTabunganCount
		}

		//END OF TRANS SETORAN TABUNGAN

		//TRANS TARIK TABUNGAN
		sqlStatement = "SELECT coalesce(SUM(pokok),0) principal, COUNT(t.no_rekening) count " +
			"FROM " +
			"tabtrans t," +
			"nasabah n," +
			"tabung m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND m.kode_produk not in ('ABP','AKP','KPE') " +
			"AND (t.tgl_trans) = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=2 and pokok>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.TransTarikTabungan, &messageList.TransTarikTabunganCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			TransTarikTabungan = messageList.TransTarikTabungan
			TransTarikTabunganCount = messageList.TransTarikTabunganCount
		}

		//END OF TRANS TARIK TABUNGAN

		//TRANS SETORAN DEPOSITO
		sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) principal, COUNT(t.no_rekening) count " +
			"FROM " +
			"deptrans t," +
			"nasabah n," +
			"deposito m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND (t.tgl_trans) = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=1 and pokok_trans>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.TransSetorDeposito, &messageList.TransSetorDepositoCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			TransSetorDeposito = messageList.TransSetorDeposito
			TransSetorDepositoCount = messageList.TransSetorDepositoCount
		}

		//END OF TRANS SETORAN DEPOSITO

		//TRANS TARIK DEPOSITO
		sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) principal, COUNT(t.no_rekening) count " +
			"FROM " +
			"deptrans t," +
			"nasabah n," +
			"deposito m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND (t.tgl_trans) = '" + cTgl + "' " +
			"and floor(my_kode_trans/100)=2 and pokok_trans>0 and m.kode_kantor='" + KodeKantor + "' "
		rows, err = db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReal{}
			err = rows.Scan(&messageList.TransTarikDeposito, &messageList.TransTarikDepositoCount)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			TransTarikDeposito = messageList.TransTarikDeposito
			TransTarikDepositoCount = messageList.TransTarikDepositoCount
		}

		//END OF TRANS TARIK DEPOSITO

		///DELETE DULU DATA
		SendAPIDelete("delDashboardReal" + "/" + KodeKantor)
		//INSERT INTO DASHBOARD_REAL
		t := time.Now()
		ts := t.Format("2006-01-02 15:04:05")
		//ts := t.Format("2006-01-02")

		a := fmt.Sprintf("%.2f", RealisasiToday)
		b := strconv.Itoa(RealisasiTodayCount)
		c := fmt.Sprintf("%.2f", RealisasiMonth)
		d := strconv.Itoa(RealisasiMonthCount)
		e := fmt.Sprintf("%.2f", BukaTabunganToday)
		f := strconv.Itoa(BukaTabunganTodayCount)
		g := fmt.Sprintf("%.2f", BukaTabunganMonth)
		h := strconv.Itoa(BukaTabunganMonthCount)
		i := fmt.Sprintf("%.2f", BukaDepositoToday)
		j := strconv.Itoa(BukaDepositoTodayCount)
		k := fmt.Sprintf("%.2f", BukaDepositoMonth)
		l := strconv.Itoa(BukaDepositoMonthCount)
		m := fmt.Sprintf("%.2f", TransAngsuranPokok)
		n := fmt.Sprintf("%.2f", TransAngsuranBunga)
		o := fmt.Sprintf("%.2f", TransAngsuranDenda)
		p := strconv.Itoa(TransAngsuranCount)
		q := fmt.Sprintf("%.2f", TransSetorTabungan)
		r := strconv.Itoa(TransSetorTabunganCount)
		s := fmt.Sprintf("%.2f", TransTarikTabungan)
		t1 := strconv.Itoa(TransTarikTabunganCount)
		u := fmt.Sprintf("%.2f", TransSetorDeposito)
		v := strconv.Itoa(TransSetorDepositoCount)
		w := fmt.Sprintf("%.2f", TransTarikDeposito)
		x := strconv.Itoa(TransTarikDepositoCount)
		SendAPIPost("setDashboardReal" + "/" + a + "/" + b + "/" + c + "/" + d + "/" + e + "/" +
			f + "/" + g + "/" + h + "/" + i + "/" + j + "/" + k + "/" + l + "/" + m + "/" + n + "/" + o + "/" +
			p + "/" + q + "/" + r + "/" + s + "/" + t1 + "/" + u + "/" + v + "/" + w + "/" + x + "/" + ts + "/" + KodeKantor)

	}

	//REALISASI HARI INI
	sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiToday,COUNT(m.no_rekening) RealisasiTodayCount " +
		"FROM " +
		"kretrans t," +
		"nasabah n," +
		"kredit m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND m.tgl_realisasi = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=1 and pokok>0 "
	rows, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.RealisasiToday, &messageList.RealisasiTodayCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		RealisasiToday = messageList.RealisasiToday
		RealisasiTodayCount = messageList.RealisasiTodayCount

	}

	//END OF REALISASI HARI INI

	//REALISASI BULAN INI
	sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
		"FROM " +
		"kretrans t," +
		"nasabah n," +
		"kredit m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND MONTH(m.tgl_realisasi) = '" + cMonth + "' " +
		"AND YEAR(m.tgl_realisasi) = '" + cYear + "' " +
		"and floor(my_kode_trans/100)=1  and pokok>0"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.RealisasiMonth, &messageList.RealisasiMonthCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		RealisasiMonth = messageList.RealisasiMonth
		RealisasiMonthCount = messageList.RealisasiMonthCount
	}

	//END OF REALISASI BULAN INI

	//BUKA TABUNGAN HARI INI
	sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
		"FROM " +
		"tabtrans t," +
		"nasabah n," +
		"tabung m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND m.kode_produk not in ('ABP','AKP','KPE') " +
		"AND (m.tgl_register) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=1 and pokok>0 "
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.BukaTabunganToday, &messageList.BukaTabunganTodayCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		BukaTabunganToday = messageList.BukaTabunganToday
		BukaTabunganTodayCount = messageList.BukaTabunganTodayCount
	}

	//END OF BUKA TABUNGAN HARI INI

	//BUKA TABUNGAN BULAN INI
	sqlStatement = "SELECT coalesce(SUM(pokok),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
		"FROM " +
		"tabtrans t," +
		"nasabah n," +
		"tabung m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND MONTH(m.tgl_register) = '" + cMonth + "' " +
		"AND YEAR(m.tgl_register) = '" + cYear + "' " +
		"and floor(my_kode_trans/100)=1 and pokok>0 "
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.BukaTabunganMonth, &messageList.BukaTabunganMonthCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		BukaTabunganMonth = messageList.BukaTabunganMonth
		BukaTabunganMonthCount = messageList.BukaTabunganMonthCount
	}

	//END OF BUKA TABUNGAN BULAN INI

	//BUKA DEPOSITO HARI INI
	sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
		"FROM " +
		"deptrans t," +
		"nasabah n," +
		"deposito m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND (m.tgl_registrasi) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=1 and pokok_trans>0"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.BukaDepositoToday, &messageList.BukaDepositoTodayCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		BukaDepositoToday = messageList.BukaDepositoToday
		BukaDepositoTodayCount = messageList.BukaDepositoTodayCount
	}

	//END OF BUKA DEPOSITO HARI INI

	//BUKA DEPOSITO BULAN INI
	sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) RealisasiMonth,COUNT(m.no_rekening) RealisasiMonthCount " +
		"FROM " +
		"deptrans t," +
		"nasabah n," +
		"deposito m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND MONTH(m.tgl_registrasi) = '" + cMonth + "' " +
		"AND YEAR(m.tgl_registrasi) = '" + cYear + "' " +
		"and floor(my_kode_trans/100)=1 and pokok_trans>0"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.BukaDepositoMonth, &messageList.BukaDepositoMonthCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		BukaDepositoMonth = messageList.BukaDepositoMonth
		BukaDepositoMonthCount = messageList.BukaDepositoMonthCount
	}

	//END OF BUKA DEPOSITO BULAN INI

	//TRANS ANGSURAN
	sqlStatement = "SELECT coalesce(SUM(pokok),0) principal, coalesce(SUM(bunga),0) interest, coalesce(SUM(denda),0) fines, COUNT(t.no_rekening) count " +
		"FROM " +
		"kretrans t," +
		"nasabah n," +
		"kredit m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND (t.tgl_trans) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=3 and (pokok<>0 or bunga<>0 or denda<>0)"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.TransAngsuranPokok, &messageList.TransAngsuranBunga,
			&messageList.TransAngsuranDenda, &messageList.TransAngsuranCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		TransAngsuranPokok = messageList.TransAngsuranPokok
		TransAngsuranBunga = messageList.TransAngsuranBunga
		TransAngsuranDenda = messageList.TransAngsuranDenda
		TransAngsuranCount = messageList.TransAngsuranCount
	}

	//END OF TRANS ANGSURAN

	//TRANS SETORAN TABUNGAN
	sqlStatement = "SELECT coalesce(SUM(pokok),0) principal, COUNT(t.no_rekening) count " +
		"FROM " +
		"tabtrans t," +
		"nasabah n," +
		"tabung m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND m.kode_produk not in ('ABP','AKP','KPE') " +
		"AND (t.tgl_trans) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=1 and pokok>0"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.TransSetorTabungan, &messageList.TransSetorTabunganCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		TransSetorTabungan = messageList.TransSetorTabungan
		TransSetorTabunganCount = messageList.TransSetorTabunganCount
	}

	//END OF TRANS SETORAN TABUNGAN

	//TRANS TARIK TABUNGAN
	sqlStatement = "SELECT coalesce(SUM(pokok),0) principal, COUNT(t.no_rekening) count " +
		"FROM " +
		"tabtrans t," +
		"nasabah n," +
		"tabung m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND m.kode_produk not in ('ABP','AKP','KPE') " +
		"AND (t.tgl_trans) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=2 and pokok>0"
	rows, err = db.Query(sqlStatement)

	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.TransTarikTabungan, &messageList.TransTarikTabunganCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		TransTarikTabungan = messageList.TransTarikTabungan
		TransTarikTabunganCount = messageList.TransTarikTabunganCount
	}

	//END OF TRANS TARIK TABUNGAN

	//TRANS SETORAN DEPOSITO
	sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) principal, COUNT(t.no_rekening) count " +
		"FROM " +
		"deptrans t," +
		"nasabah n," +
		"deposito m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND (t.tgl_trans) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=1 and pokok_trans>0"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.TransSetorDeposito, &messageList.TransSetorDepositoCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		TransSetorDeposito = messageList.TransSetorDeposito
		TransSetorDepositoCount = messageList.TransSetorDepositoCount
	}

	//END OF TRANS SETORAN DEPOSITO

	//TRANS TARIK DEPOSITO
	sqlStatement = "SELECT coalesce(SUM(pokok_trans),0) principal, COUNT(t.no_rekening) count " +
		"FROM " +
		"deptrans t," +
		"nasabah n," +
		"deposito m " +
		"WHERE " +
		"t.no_rekening = m.no_rekening " +
		"AND n.nasabah_id = m.nasabah_id " +
		"AND (t.tgl_trans) = '" + cTgl + "' " +
		"and floor(my_kode_trans/100)=2 and pokok_trans>0"
	rows, err = db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReal{}
		err = rows.Scan(&messageList.TransTarikDeposito, &messageList.TransTarikDepositoCount)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		TransTarikDeposito = messageList.TransTarikDeposito
		TransTarikDepositoCount = messageList.TransTarikDepositoCount
	}

	//END OF TRANS TARIK DEPOSITO

	///DELETE DULU DATA

	SendAPIDelete("delDashboardReal" + "/All")
	//INSERT INTO DASHBOARD_REAL
	t := time.Now()
	ts := t.Format("2006-01-02 15:04:05")

	a := fmt.Sprintf("%.2f", RealisasiToday)
	b := strconv.Itoa(RealisasiTodayCount)
	c := fmt.Sprintf("%.2f", RealisasiMonth)
	d := strconv.Itoa(RealisasiMonthCount)
	e := fmt.Sprintf("%.2f", BukaTabunganToday)
	f := strconv.Itoa(BukaTabunganTodayCount)
	g := fmt.Sprintf("%.2f", BukaTabunganMonth)
	h := strconv.Itoa(BukaTabunganMonthCount)
	i := fmt.Sprintf("%.2f", BukaDepositoToday)
	j := strconv.Itoa(BukaDepositoTodayCount)
	k := fmt.Sprintf("%.2f", BukaDepositoMonth)
	l := strconv.Itoa(BukaDepositoMonthCount)
	m := fmt.Sprintf("%.2f", TransAngsuranPokok)
	n := fmt.Sprintf("%.2f", TransAngsuranBunga)
	o := fmt.Sprintf("%.2f", TransAngsuranDenda)
	p := strconv.Itoa(TransAngsuranCount)
	q := fmt.Sprintf("%.2f", TransSetorTabungan)
	r := strconv.Itoa(TransSetorTabunganCount)
	s := fmt.Sprintf("%.2f", TransTarikTabungan)
	t1 := strconv.Itoa(TransTarikTabunganCount)
	u := fmt.Sprintf("%.2f", TransSetorDeposito)
	v := strconv.Itoa(TransSetorDepositoCount)
	w := fmt.Sprintf("%.2f", TransTarikDeposito)
	x := strconv.Itoa(TransTarikDepositoCount)
	SendAPIPost("setDashboardReal" + "/" + a + "/" + b + "/" + c + "/" + d + "/" + e + "/" +
		f + "/" + g + "/" + h + "/" + i + "/" + j + "/" + k + "/" + l + "/" + m + "/" + n + "/" + o + "/" +
		p + "/" + q + "/" + r + "/" + s + "/" + t1 + "/" + u + "/" + v + "/" + w + "/" + x + "/" + ts + "/All")

	defer db.Close()
	functions.Logger().Info("Successfully Get Data Dashboard Real ")
}

func GetDataDashboardChart() {
	functions.Logger().Info("Starting Scheduler Get Data Dashboard Chart ")
	//cTgl := jodaTime.Format("YYYYMMdd", time.Now())
	//cMonth := jodaTime.Format("MM", time.Now())

	cYear := jodaTime.Format("YYYY", time.Now())
	DBRe := os.Getenv("DB_DATABASE_RE")

	sqlStatement := "SELECT kode_kantor from app_kode_kantor order by kode_kantor "
	rowsY, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}
	for rowsY.Next() {
		messageList := SelectDashboardReal{}
		err = rowsY.Scan(&messageList.KodeKantor)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		KodeKantor := messageList.KodeKantor
		//CHART PEMBUKAAN TABUNGAN
		///DELETE DULU DATA
		SendAPIDelete("delOpenTabungan" + "/" + KodeKantor)
		///DELETE DULU VIEW

		m, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		for i := 1; i <= m; i++ {
			s := strconv.Itoa(i)

			_, err := db.Query("drop view if exists cdashboard_chart_open_tabungan ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			_, err = db.Query("create view cdashboard_chart_open_tabungan " +
				"as " +
				"SELECT t.no_rekening,pokok " +
				"from tabtrans t,tabung m,nasabah n " +
				"where t.no_rekening=m.no_rekening " +
				"and n.nasabah_id=m.nasabah_id " +
				"and floor(my_kode_trans/100)=1 " +
				"AND m.kode_produk not in ('ABP','AKP','KPE') " +
				"and t.tgl_trans=m.tgl_register " +
				"and year(m.tgl_register)= '" + cYear + "' " +
				"and month(m.tgl_register)='" + s + "' and m.kode_kantor='" + KodeKantor + "' ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_tabungan "
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)

				SendAPIPost("setOpenTabungan" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}

		//END OF CHART PEMBUKAAN TABUNGAN

		//CHART PEMBUKAAN DEPOSITO
		///DELETE DULU DATA
		SendAPIDelete("delOpenDeposito" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		for i := 1; i <= m; i++ {
			s := strconv.Itoa(i)

			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_open_deposito ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			_, err = db.Query("create view cdashboard_chart_open_deposito " +
				"as " +
				"SELECT t.no_rekening,pokok_trans pokok " +
				"from deptrans t,deposito m,nasabah n " +
				"where t.no_rekening=m.no_rekening " +
				"and n.nasabah_id=m.nasabah_id " +
				"and floor(my_kode_trans/100)=1 " +
				"and t.tgl_trans=m.tgl_registrasi " +
				"and year(m.tgl_registrasi)= '" + cYear + "' " +
				"and month(m.tgl_registrasi)='" + s + "' and m.kode_kantor='" + KodeKantor + "' ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_deposito "
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setOpenDeposito" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}
		//END OF CHART PEMBUKAAN DEPOSITO

		//CHART PEMBUKAAN KREDIT
		///DELETE DULU DATA
		SendAPIDelete("delOpenKredit" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		for i := 1; i <= m; i++ {
			s := strconv.Itoa(i)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_penyaluran_kredit ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			_, err = db.Query("create view cdashboard_chart_penyaluran_kredit " +
				"as " +
				"SELECT t.no_rekening,pokok " +
				"from kretrans t,kredit m,nasabah n " +
				"where t.no_rekening=m.no_rekening " +
				"and n.nasabah_id=m.nasabah_id " +
				"and floor(my_kode_trans/100)=1 " +
				"and t.tgl_trans=m.tgl_realisasi " +
				"and year(m.tgl_realisasi)= '" + cYear + "' " +
				"and month(m.tgl_realisasi)='" + s + "' and m.kode_kantor='" + KodeKantor + "' ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_penyaluran_kredit "
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setOpenKredit" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}
		//END OF CHART PEMBUKAAN KREDIT

		//CHART SALDO TABUNGAN
		///DELETE DULU DATA
		SendAPIDelete("delSaldoTabungan" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		dd, _ := strconv.Atoi(jodaTime.Format("dd", time.Now()))
		ds := strconv.Itoa(dd)
		if dd < 10 {
			ds = "0" + ds
		}
		for i := 1; i <= m; i++ {
			z := i
			s := strconv.Itoa(i)
			d := ""
			if z < 10 {
				d = "0" + s
			} else {
				d = s
			}
			newdate := cYear + "-" + d + "-01"
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 1, -1)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				newdate := cYear + "-" + d + "-" + ds
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_tabungan ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".tabung_" + cNewDate + " a,tabung where tabung.no_Rekening=a.no_rekening and tabung.kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {
				_, err = db.Query("create view cdashboard_chart_saldo_tabungan " +
					"as " +
					"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".tabung_" + cNewDate + " a,tabung " +
					"where tabung.no_rekening=a.no_rekening and tabung.kode_kantor='" + KodeKantor + "' " +
					"AND tabung.kode_produk not in ('ABP','AKP','KPE') ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_tabungan"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Count)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					z := strconv.Itoa(i)
					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Count)
					SendAPIPost("setSaldoTabungan" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART SALDO TABUNGAN

		//CHART SALDO DEPOSITO
		///DELETE DULU DATA

		SendAPIDelete("delSaldoDeposito" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
		ds = strconv.Itoa(dd)
		if dd < 10 {
			ds = "0" + ds
		}
		for i := 1; i <= m; i++ {
			z := i
			s := strconv.Itoa(i)
			d := ""
			if z < 10 {
				d = "0" + s
			} else {
				d = s
			}
			newdate := cYear + "-" + d + "-01"
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 1, -1)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				newdate := cYear + "-" + d + "-" + ds
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}

			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_deposito ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".deposito_" + cNewDate + " a,deposito" +
				" where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {

				_, err := db.Query("create view cdashboard_chart_saldo_deposito " +
					"as " +
					"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".deposito_" + cNewDate + " a,deposito " +
					"where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_deposito"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Count)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					z := strconv.Itoa(i)
					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Count)
					SendAPIPost("setSaldoDeposito" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART SALDO DEPOSITO

		//CHART SALDO KREDIT

		///DELETE DULU DATA

		SendAPIDelete("delSaldoKredit" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
		ds = strconv.Itoa(dd)
		if dd < 10 {
			ds = "0" + ds
		}
		for i := 1; i <= m; i++ {
			z := i
			s := strconv.Itoa(i)
			d := ""
			if z < 10 {
				d = "0" + s
			} else {
				d = s
			}
			newdate := cYear + "-" + d + "-01"
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 1, -1)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				newdate := cYear + "-" + d + "-" + ds
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_kredit ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a, kredit " +
				"where kredit.no_rekening=a.no_rekening and kredit.kode_kantor='" + KodeKantor + "'  limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {
				_, err := db.Query("create view cdashboard_chart_saldo_kredit " +
					"as " +
					"SELECT a.no_rekening,a.baki_debet amount from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening " +
					"and kredit.kode_kantor='" + KodeKantor + "'  ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_kredit"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Count)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					z := strconv.Itoa(i)
					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Count)
					SendAPIPost("setSaldoKredit" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART SALDO KREDIT

		//CHART PERTUMBUHAN ANGSURAN

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganAngsuran" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
		ds = strconv.Itoa(dd)
		if dd < 10 {
			ds = "0" + ds
		}
		for i := 1; i <= m; i++ {
			z := i
			s := strconv.Itoa(i)
			d := ""
			if z < 10 {
				d = "0" + s
			} else {
				d = s
			}
			newdate := cYear + "-" + d + "-01"
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 1, -1)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				newdate := cYear + "-" + d + "-" + ds
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}

			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_angsuran ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			_, err = db.Query("create view cdashboard_chart_pertumbuhan_angsuran " +
				"as " +
				"SELECT t.no_rekening, " +
				"coalesce(SUM(if(floor(my_kode_trans/100)=3,pokok,0)),0) amount," +
				"coalesce(SUM(if(floor(my_kode_trans/100)=3,bunga,0)),0) amount_bunga " +
				"FROM " +
				"kretrans t," +
				"nasabah n," +
				"kredit m " +
				"WHERE " +
				"t.no_rekening = m.no_rekening " +
				"AND n.nasabah_id = m.nasabah_id " +
				"AND tgl_trans>=date_add(date_add(LAST_DAY(last_day('" + cNewDate + "')),interval 1 DAY),interval -1 MONTH)  " +
				"AND tgl_trans<=last_day('" + cNewDate + "') and m.kode_kantor='" + KodeKantor + "' group by t.no_rekening having amount<>0 or amount_bunga<>0 ")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,coalesce(sum(amount_bunga),0) amount_bunga from cdashboard_chart_pertumbuhan_angsuran"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Amount2)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Amount2)
				SendAPIPost("setPerkembanganAngsuran" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}
		//END OF CHART PERTUMBUHAN ANGSURAN

		//CHART PERKEMBANGAN NPL

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganNPL" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
		ds = strconv.Itoa(dd)
		if dd < 10 {
			ds = "0" + ds
		}
		for i := 1; i <= m; i++ {
			z := i
			s := strconv.Itoa(i)
			d := ""
			if z < 10 {
				d = "0" + s
			} else {
				d = s
			}
			newdate := cYear + "-" + d + "-01"
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 1, -1)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				newdate := cYear + "-" + d + "-" + ds
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_perkembangan_npl ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {

				_, err := db.Query("create view cdashboard_chart_perkembangan_npl " +
					"as " +
					"SELECT coalesce(if(sum(a.baki_debet)>0," +
					"sum(if(a.my_kolek_number>2,a.baki_debet,0))/sum(a.baki_debet)*100,0),0) as npl from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
					"where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "'")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(npl,0) npl from cdashboard_chart_perkembangan_npl"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}
					z := strconv.Itoa(i)
					amount := fmt.Sprintf("%.2f", messageList.Amount)
					SendAPIPost("setPerkembanganNPL" + "/" + z + "/" + amount + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART PERKEMBANGAN NPL

		//CHART PERKEMBANGAN ASET,KEWAJIBAN,MODAL,LABARUGI

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganAset" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
		ds = strconv.Itoa(dd)
		if dd < 10 {
			ds = "0" + ds
		}
		for i := 1; i <= m; i++ {
			z := i
			s := strconv.Itoa(i)
			d := ""
			if z < 10 {
				d = "0" + s
			} else {
				d = s
			}
			newdate := cYear + "-" + d + "-01"
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 1, -1)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				newdate := cYear + "-" + d + "-" + ds
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_aset ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {

				_, err := db.Query("create view cdashboard_chart_pertumbuhan_aset " +
					"as " +
					"SELECT sum(if(kode_perk='1',saldo_akhir,0)) aset,sum(if(kode_perk='2',saldo_akhir,0)) kewajiban," +
					"sum(if(kode_perk='3',saldo_akhir,0)) modal,sum(if(kode_perk='4',saldo_akhir,0)) pendapatan," +
					"sum(if(kode_perk='5',saldo_akhir,0)) biaya " +
					"from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(aset,0) aset," +
					" coalesce(kewajiban,0) kewajiban," +
					"coalesce(modal,0) modal,coalesce(pendapatan,0) pendapatan," +
					"coalesce(biaya,0) biaya from cdashboard_chart_pertumbuhan_aset"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Amount2, &messageList.Amount3, &messageList.Amount4, &messageList.Amount5)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					LabaRugi := messageList.Amount4 - messageList.Amount5
					z := strconv.Itoa(i)
					aset := fmt.Sprintf("%.2f", messageList.Amount)
					kewajiban := fmt.Sprintf("%.2f", messageList.Amount2)
					modal := fmt.Sprintf("%.2f", messageList.Amount3)
					laba := fmt.Sprintf("%.2f", LabaRugi)
					pendapatan := fmt.Sprintf("%.2f", messageList.Amount4)
					biaya := fmt.Sprintf("%.2f", messageList.Amount5)
					SendAPIPost("setPerkembanganAset" + "/" + z + "/" + aset + "/" + kewajiban + "/" + modal + "/" + laba + "/" + pendapatan + "/" + biaya + "/" + KodeKantor)

				}
			}

		}
		//END OF CHART PERKEMBANGAN ASET
	}

	//CHART PEMBUKAAN TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delOpenTabungan" + "/" + "All")

	m, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
	for i := 1; i <= m; i++ {
		s := strconv.Itoa(i)

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_open_tabungan ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_open_tabungan " +
			"as " +
			"SELECT t.no_rekening,pokok " +
			"from tabtrans t,tabung m,nasabah n " +
			"where t.no_rekening=m.no_rekening " +
			"and n.nasabah_id=m.nasabah_id " +
			"AND m.kode_produk not in ('ABP','AKP','KPE') " +
			"and floor(my_kode_trans/100)=1 " +
			"and t.tgl_trans=m.tgl_register " +
			"and year(m.tgl_register)= '" + cYear + "' " +
			"and month(m.tgl_register)='" + s + "'")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_tabungan "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)

			SendAPIPost("setOpenTabungan" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}

	//END OF CHART PEMBUKAAN TABUNGAN

	//CHART PEMBUKAAN DEPOSITO
	///DELETE DULU DATA
	SendAPIDelete("delOpenDeposito" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	for i := 1; i <= m; i++ {
		s := strconv.Itoa(i)

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_open_deposito ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_open_deposito " +
			"as " +
			"SELECT t.no_rekening,pokok_trans pokok " +
			"from deptrans t,deposito m,nasabah n " +
			"where t.no_rekening=m.no_rekening " +
			"and n.nasabah_id=m.nasabah_id " +
			"and floor(my_kode_trans/100)=1 " +
			"and t.tgl_trans=m.tgl_registrasi " +
			"and year(m.tgl_registrasi)= '" + cYear + "' " +
			"and month(m.tgl_registrasi)='" + s + "'")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_deposito "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			SendAPIPost("setOpenDeposito" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}
	//END OF CHART PEMBUKAAN DEPOSITO

	//CHART PEMBUKAAN KREDIT
	///DELETE DULU DATA
	SendAPIDelete("delOpenKredit" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	for i := 1; i <= m; i++ {
		s := strconv.Itoa(i)

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_penyaluran_kredit ")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_penyaluran_kredit " +
			"as " +
			"SELECT t.no_rekening,pokok " +
			"from kretrans t,kredit m,nasabah n " +
			"where t.no_rekening=m.no_rekening " +
			"and n.nasabah_id=m.nasabah_id " +
			"and floor(my_kode_trans/100)=1 " +
			"and t.tgl_trans=m.tgl_realisasi " +
			"and year(m.tgl_realisasi)= '" + cYear + "' " +
			"and month(m.tgl_realisasi)='" + s + "'")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_penyaluran_kredit "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			SendAPIPost("setOpenKredit" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}

	//END OF CHART PEMBUKAAN KREDIT

	//CHART SALDO TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delSaldoTabungan" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	dd, _ := strconv.Atoi(jodaTime.Format("dd", time.Now()))
	ds := strconv.Itoa(dd)
	if dd < 10 {
		ds = "0" + ds
	}
	for i := 1; i <= m; i++ {
		z := i
		s := strconv.Itoa(i)
		d := ""
		if z < 10 {
			d = "0" + s
		} else {
			d = s
		}
		newdate := cYear + "-" + d + "-01"
		t, _ := time.Parse("2006-01-02", newdate)
		t = t.AddDate(0, 1, -1)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			newdate := cYear + "-" + d + "-" + ds
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_tabungan ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".tabung_" + cNewDate + " limit 1"
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_saldo_tabungan " +
				"as " +
				"SELECT no_rekening,saldo_akhir amount from " + DBRe + ".tabung_" + cNewDate + " " +
				"where no_rekening not in (select no_rekening from tabung where kode_produk in ('ABP','AKP','KPE'))  ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_tabungan"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setSaldoTabungan" + "/" + z + "/" + amount + "/" + count + "/All")

			}

		}

	}
	//END OF CHART SALDO TABUNGAN

	//CHART SALDO DEPOSITO
	///DELETE DULU DATA

	SendAPIDelete("delSaldoDeposito" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
	ds = strconv.Itoa(dd)
	if dd < 10 {
		ds = "0" + ds
	}
	for i := 1; i <= m; i++ {
		z := i
		s := strconv.Itoa(i)
		d := ""
		if z < 10 {
			d = "0" + s
		} else {
			d = s
		}
		newdate := cYear + "-" + d + "-01"
		t, _ := time.Parse("2006-01-02", newdate)
		t = t.AddDate(0, 1, -1)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			newdate := cYear + "-" + d + "-" + ds
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_deposito ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".deposito_" + cNewDate + " limit 1"
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_saldo_deposito " +
				"as " +
				"SELECT no_rekening,saldo_akhir amount from " + DBRe + ".deposito_" + cNewDate + " ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_deposito"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setSaldoDeposito" + "/" + z + "/" + amount + "/" + count + "/All")

			}

		}

	}
	//END OF CHART SALDO DEPOSITO

	//CHART SALDO KREDIT

	///DELETE DULU DATA

	SendAPIDelete("delSaldoKredit" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
	ds = strconv.Itoa(dd)
	if dd < 10 {
		ds = "0" + ds
	}
	for i := 1; i <= m; i++ {
		z := i
		s := strconv.Itoa(i)
		d := ""
		if z < 10 {
			d = "0" + s
		} else {
			d = s
		}
		newdate := cYear + "-" + d + "-01"
		t, _ := time.Parse("2006-01-02", newdate)
		t = t.AddDate(0, 1, -1)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			newdate := cYear + "-" + d + "-" + ds
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_kredit ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where kredit.no_rekening = a.no_rekening limit 1 "
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_saldo_kredit " +
				"as " +
				"SELECT a.no_rekening,a.baki_debet amount from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
				"where a.no_rekening=kredit.no_rekening ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_kredit"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setSaldoKredit" + "/" + z + "/" + amount + "/" + count + "/All")

			}

		}

	}
	//END OF CHART SALDO KREDIT

	//CHART PERTUMBUHAN ANGSURAN

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganAngsuran" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
	ds = strconv.Itoa(dd)
	if dd < 10 {
		ds = "0" + ds
	}
	for i := 1; i <= m; i++ {
		z := i
		s := strconv.Itoa(i)
		d := ""
		if z < 10 {
			d = "0" + s
		} else {
			d = s
		}
		newdate := cYear + "-" + d + "-01"
		t, _ := time.Parse("2006-01-02", newdate)
		t = t.AddDate(0, 1, -1)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			newdate := cYear + "-" + d + "-" + ds
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_angsuran ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_pertumbuhan_angsuran " +
			"as " +
			"SELECT t.no_rekening, " +
			"coalesce(SUM(if(floor(my_kode_trans/100)=3,pokok,0)),0) amount," +
			"coalesce(SUM(if(floor(my_kode_trans/100)=3,bunga,0)),0) amount_bunga " +
			"FROM " +
			"kretrans t," +
			"nasabah n," +
			"kredit m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND tgl_trans>=date_add(date_add(LAST_DAY(last_day('" + cNewDate + "')),interval 1 DAY),interval -1 MONTH)  " +
			"AND tgl_trans<=last_day('" + cNewDate + "') group by t.no_rekening having amount<>0 or amount_bunga<>0 ")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select coalesce(sum(amount),0) amount,coalesce(sum(amount_bunga),0) amount_bunga from cdashboard_chart_pertumbuhan_angsuran"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Amount2)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Amount2)
			SendAPIPost("setPerkembanganAngsuran" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}
	//END OF CHART PERTUMBUHAN ANGSURAN

	//CHART PERKEMBANGAN NPL

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganNPL" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
	ds = strconv.Itoa(dd)
	if dd < 10 {
		ds = "0" + ds
	}
	for i := 1; i <= m; i++ {
		z := i
		s := strconv.Itoa(i)
		d := ""
		if z < 10 {
			d = "0" + s
		} else {
			d = s
		}
		newdate := cYear + "-" + d + "-01"
		t, _ := time.Parse("2006-01-02", newdate)
		t = t.AddDate(0, 1, -1)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			newdate := cYear + "-" + d + "-" + ds
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_perkembangan_npl ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening limit 1 "
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_perkembangan_npl " +
				"as " +
				"SELECT coalesce(if(sum(a.baki_debet)>0," +
				"sum(if(a.my_kolek_number>2,a.baki_debet,0))/sum(a.baki_debet)*100,0),0) as npl from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
				"where  a.no_rekening=kredit.no_rekening ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(npl,0) npl from cdashboard_chart_perkembangan_npl"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				SendAPIPost("setPerkembanganNPL" + "/" + z + "/" + amount + "/All")

			}

		}

	}
	//END OF CHART PERKEMBANGAN NPL

	//CHART PERKEMBANGAN ASET,KEWAJIBAN,MODAL,LABARUGI

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganAset" + "/All")

	m, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
	dd, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
	ds = strconv.Itoa(dd)
	if dd < 10 {
		ds = "0" + ds
	}
	for i := 1; i <= m; i++ {
		z := i
		s := strconv.Itoa(i)
		d := ""
		if z < 10 {
			d = "0" + s
		} else {
			d = s
		}
		newdate := cYear + "-" + d + "-01"
		t, _ := time.Parse("2006-01-02", newdate)
		t = t.AddDate(0, 1, -1)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			newdate := cYear + "-" + d + "-" + ds
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_aset ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1 "
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_chart_pertumbuhan_aset " +
				"as " +
				"SELECT sum(if(kode_perk='1',saldo_akhir,0)) aset,sum(if(kode_perk='2',saldo_akhir,0)) kewajiban," +
				"sum(if(kode_perk='3',saldo_akhir,0)) modal,sum(if(kode_perk='4',saldo_akhir,0)) pendapatan," +
				"sum(if(kode_perk='5',saldo_akhir,0)) biaya " +
				"from " + DBRe + ".acc_" + cNewDate + " ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(aset,0) aset," +
				" coalesce(kewajiban,0) kewajiban," +
				"coalesce(modal,0) modal,coalesce(pendapatan,0) pendapatan," +
				"coalesce(biaya,0) biaya from cdashboard_chart_pertumbuhan_aset"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Amount2, &messageList.Amount3, &messageList.Amount4, &messageList.Amount5)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				LabaRugi := messageList.Amount4 - messageList.Amount5
				z := strconv.Itoa(i)
				aset := fmt.Sprintf("%.2f", messageList.Amount)
				kewajiban := fmt.Sprintf("%.2f", messageList.Amount2)
				modal := fmt.Sprintf("%.2f", messageList.Amount3)
				laba := fmt.Sprintf("%.2f", LabaRugi)
				pendapatan := fmt.Sprintf("%.2f", messageList.Amount4)
				biaya := fmt.Sprintf("%.2f", messageList.Amount5)
				SendAPIPost("setPerkembanganAset" + "/" + z + "/" + aset + "/" + kewajiban + "/" + modal + "/" + laba + "/" + pendapatan + "/" + biaya + "/All")

			}

		}

	}
	//END OF CHART PERKEMBANGAN ASET
	defer db.Close()
	functions.Logger().Info("Successfully Get Data Dashboard Chart ")
}

func GetDataDashboardChartTahun() {
	functions.Logger().Info("Starting Scheduler Get Data Dashboard Chart Tahunan ")

	//cYear := jodaTime.Format("YYYY", time.Now())
	DBRe := os.Getenv("DB_DATABASE_RE")
	m, _ := strconv.Atoi(jodaTime.Format("YYYY", time.Now()))

	sqlStatement := "SELECT kode_kantor from app_kode_kantor order by kode_kantor "
	rowsY, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsY.Next() {
		messageList := SelectDashboardReal{}
		err = rowsY.Scan(&messageList.KodeKantor)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		KodeKantor := messageList.KodeKantor

		//CHART PEMBUKAAN TABUNGAN
		///DELETE DULU DATA
		SendAPIDelete("delOpenTabunganTahun" + "/" + KodeKantor)
		///DELETE DULU VIEW

		m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)

			_, err := db.Query("drop view if exists cdashboard_chart_open_tabungan_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			_, err = db.Query("create view cdashboard_chart_open_tabungan_tahun " +
				"as " +
				"SELECT t.no_rekening,pokok " +
				"from tabtrans t,tabung m,nasabah n " +
				"where t.no_rekening=m.no_rekening " +
				"and n.nasabah_id=m.nasabah_id " +
				"and floor(my_kode_trans/100)=1 " +
				"AND m.kode_produk not in ('ABP','AKP','KPE') " +
				"and t.tgl_trans=m.tgl_register " +
				"and year(m.tgl_register)= '" + s + "' " +
				"and m.kode_kantor='" + KodeKantor + "' ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_tabungan_tahun "
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)

				SendAPIPost("setOpenTabunganTahun" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}

		//END OF CHART PEMBUKAAN TABUNGAN

		//CHART PEMBUKAAN DEPOSITO
		///DELETE DULU DATA
		SendAPIDelete("delOpenDepositoTahun" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)

			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_open_deposito_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			_, err = db.Query("create view cdashboard_chart_open_deposito_tahun " +
				"as " +
				"SELECT t.no_rekening,pokok_trans pokok " +
				"from deptrans t,deposito m,nasabah n " +
				"where t.no_rekening=m.no_rekening " +
				"and n.nasabah_id=m.nasabah_id " +
				"and floor(my_kode_trans/100)=1 " +
				"and t.tgl_trans=m.tgl_registrasi " +
				"and year(m.tgl_registrasi)= '" + s + "' " +
				"and m.kode_kantor='" + KodeKantor + "' ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_deposito_tahun "
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setOpenDepositoTahun" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}
		//END OF CHART PEMBUKAAN DEPOSITO

		//CHART PEMBUKAAN KREDIT
		///DELETE DULU DATA
		SendAPIDelete("delOpenKreditTahun" + "/" + KodeKantor)

		m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_penyaluran_kredit_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			_, err = db.Query("create view cdashboard_chart_penyaluran_kredit_tahun " +
				"as " +
				"SELECT t.no_rekening,pokok " +
				"from kretrans t,kredit m,nasabah n " +
				"where t.no_rekening=m.no_rekening " +
				"and n.nasabah_id=m.nasabah_id " +
				"and floor(my_kode_trans/100)=1 " +
				"and t.tgl_trans=m.tgl_realisasi " +
				"and year(m.tgl_realisasi)= '" + s + "' " +
				"and m.kode_kantor='" + KodeKantor + "' ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_penyaluran_kredit_tahun "
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setOpenKreditTahun" + "/" + z + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}
		//END OF CHART PEMBUKAAN KREDIT

		//CHART SALDO TABUNGAN
		///DELETE DULU DATA
		SendAPIDelete("delSaldoTabunganTahun" + "/" + KodeKantor)

		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			newdate := s + "-" + "12" + "-31"
			t, _ := time.Parse("2006-01-02", newdate)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
				ds := strconv.Itoa(z)
				d := ""
				if z < 10 {
					d = "0" + ds
				} else {
					d = ds
				}
				z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
				ds = strconv.Itoa(z)
				dx := ""
				if z < 10 {
					dx = "0" + ds
				} else {
					dx = ds
				}
				newdate := s + "-" + d + "-" + dx
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_tabungan_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".tabung_" + cNewDate + " a,tabung where tabung.no_Rekening=a.no_rekening and tabung.kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {
				_, err = db.Query("create view cdashboard_chart_saldo_tabungan_tahun " +
					"as " +
					"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".tabung_" + cNewDate + " a,tabung " +
					"where tabung.no_rekening=a.no_rekening and tabung.kode_kantor='" + KodeKantor + "' " +
					"AND tabung.kode_produk not in ('ABP','AKP','KPE') ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_tabungan_tahun"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Count)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Count)
					SendAPIPost("setSaldoTabunganTahun" + "/" + strconv.Itoa(i) + "/" + amount + "/" + count + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART SALDO TABUNGAN

		//CHART SALDO DEPOSITO
		///DELETE DULU DATA

		SendAPIDelete("delSaldoDepositoTahun" + "/" + KodeKantor)

		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			newdate := s + "-" + "12" + "-31"
			t, _ := time.Parse("2006-01-02", newdate)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
				ds := strconv.Itoa(z)
				d := ""
				if z < 10 {
					d = "0" + ds
				} else {
					d = ds
				}
				z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
				ds = strconv.Itoa(z)
				dx := ""
				if z < 10 {
					dx = "0" + ds
				} else {
					dx = ds
				}
				newdate := s + "-" + d + "-" + dx
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}

			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_deposito_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".deposito_" + cNewDate + " a,deposito" +
				" where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {

				_, err := db.Query("create view cdashboard_chart_saldo_deposito_tahun " +
					"as " +
					"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".deposito_" + cNewDate + " a,deposito " +
					"where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_deposito_tahun"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Count)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Count)
					SendAPIPost("setSaldoDepositoTahun" + "/" + strconv.Itoa(i) + "/" + amount + "/" + count + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART SALDO DEPOSITO

		//CHART SALDO KREDIT

		///DELETE DULU DATA

		SendAPIDelete("delSaldoKreditTahun" + "/" + KodeKantor)

		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			newdate := s + "-" + "12" + "-31"
			t, _ := time.Parse("2006-01-02", newdate)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
				ds := strconv.Itoa(z)
				d := ""
				if z < 10 {
					d = "0" + ds
				} else {
					d = ds
				}
				z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
				ds = strconv.Itoa(z)
				dx := ""
				if z < 10 {
					dx = "0" + ds
				} else {
					dx = ds
				}
				newdate := s + "-" + d + "-" + dx
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_kredit_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a, kredit " +
				"where kredit.no_rekening=a.no_rekening and kredit.kode_kantor='" + KodeKantor + "'  limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {
				_, err := db.Query("create view cdashboard_chart_saldo_kredit_tahun " +
					"as " +
					"SELECT a.no_rekening,a.baki_debet amount from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening " +
					"and kredit.kode_kantor='" + KodeKantor + "'  ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_kredit_tahun"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Count)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Count)
					SendAPIPost("setSaldoKreditTahun" + "/" + strconv.Itoa(i) + "/" + amount + "/" + count + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART SALDO KREDIT

		//CHART PERTUMBUHAN ANGSURAN

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganAngsuranTahun" + "/" + KodeKantor)

		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			newdate := s + "-" + "12" + "-31"
			t, _ := time.Parse("2006-01-02", newdate)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
				ds := strconv.Itoa(z)
				d := ""
				if z < 10 {
					d = "0" + ds
				} else {
					d = ds
				}
				z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
				ds = strconv.Itoa(z)
				dx := ""
				if z < 10 {
					dx = "0" + ds
				} else {
					dx = ds
				}
				newdate := s + "-" + d + "-" + dx
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}

			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_angsuran_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			_, err = db.Query("create view cdashboard_chart_pertumbuhan_angsuran_tahun " +
				"as " +
				"SELECT t.no_rekening, " +
				"coalesce(SUM(if(floor(my_kode_trans/100)=3,pokok,0)),0) amount," +
				"coalesce(SUM(if(floor(my_kode_trans/100)=3,bunga,0)),0) amount_bunga " +
				"FROM " +
				"kretrans t," +
				"nasabah n," +
				"kredit m " +
				"WHERE " +
				"t.no_rekening = m.no_rekening " +
				"AND n.nasabah_id = m.nasabah_id " +
				"AND tgl_trans>=date_add(date_add(LAST_DAY(last_day('" + cNewDate + "')),interval 1 DAY),interval -1 MONTH)  " +
				"AND tgl_trans<=last_day('" + cNewDate + "') and m.kode_kantor='" + KodeKantor + "' group by t.no_rekening having amount<>0 or amount_bunga<>0 ")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,coalesce(sum(amount_bunga),0) amount_bunga from cdashboard_chart_pertumbuhan_angsuran_tahun"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Amount2)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Amount2)
				SendAPIPost("setPerkembanganAngsuranTahun" + "/" + strconv.Itoa(i) + "/" + amount + "/" + count + "/" + KodeKantor)

			}
		}
		//END OF CHART PERTUMBUHAN ANGSURAN

		//CHART PERKEMBANGAN NPL

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganNPLTahun" + "/" + KodeKantor)

		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			newdate := s + "-" + "12" + "-31"
			t, _ := time.Parse("2006-01-02", newdate)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
				ds := strconv.Itoa(z)
				d := ""
				if z < 10 {
					d = "0" + ds
				} else {
					d = ds
				}
				z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
				ds = strconv.Itoa(z)
				dx := ""
				if z < 10 {
					dx = "0" + ds
				} else {
					dx = ds
				}
				newdate := s + "-" + d + "-" + dx
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_perkembangan_npl_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {

				_, err := db.Query("create view cdashboard_chart_perkembangan_npl_tahun " +
					"as " +
					"SELECT coalesce(if(sum(a.baki_debet)>0," +
					"sum(if(a.my_kolek_number>2,a.baki_debet,0))/sum(a.baki_debet)*100,0),0) as npl from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
					"where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "'")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(npl,0) npl from cdashboard_chart_perkembangan_npl_tahun"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}
					amount := fmt.Sprintf("%.2f", messageList.Amount)
					SendAPIPost("setPerkembanganNPLTahun" + "/" + strconv.Itoa(i) + "/" + amount + "/" + KodeKantor)

				}
			}
		}
		//END OF CHART PERKEMBANGAN NPL

		//CHART PERKEMBANGAN ASET,KEWAJIBAN,MODAL,LABARUGI

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganAsetTahun" + "/" + KodeKantor)

		for i := m - 5; i <= m; i++ {
			s := strconv.Itoa(i)
			newdate := s + "-" + "12" + "-31"
			t, _ := time.Parse("2006-01-02", newdate)
			cNewDate := jodaTime.Format("YYMMdd", t)
			if m == i {
				z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
				ds := strconv.Itoa(z)
				d := ""
				if z < 10 {
					d = "0" + ds
				} else {
					d = ds
				}
				z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
				ds = strconv.Itoa(z)
				dx := ""
				if z < 10 {
					dx = "0" + ds
				} else {
					dx = ds
				}
				newdate := s + "-" + d + "-" + dx
				t, _ := time.Parse("2006-01-02", newdate)
				t = t.AddDate(0, 0, -1)
				cNewDate = jodaTime.Format("YYMMdd", t)
			}
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_aset_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rowsX.Next() {

				_, err := db.Query("create view cdashboard_chart_pertumbuhan_aset_tahun " +
					"as " +
					"SELECT sum(if(kode_perk='1',saldo_akhir,0)) aset,sum(if(kode_perk='2',saldo_akhir,0)) kewajiban," +
					"sum(if(kode_perk='3',saldo_akhir,0)) modal,sum(if(kode_perk='4',saldo_akhir,0)) pendapatan," +
					"sum(if(kode_perk='5',saldo_akhir,0)) biaya " +
					"from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(aset,0) aset," +
					" coalesce(kewajiban,0) kewajiban," +
					"coalesce(modal,0) modal,coalesce(pendapatan,0) pendapatan," +
					"coalesce(biaya,0) biaya from cdashboard_chart_pertumbuhan_aset_tahun"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Amount2, &messageList.Amount3, &messageList.Amount4, &messageList.Amount5)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					LabaRugi := messageList.Amount4 - messageList.Amount5
					aset := fmt.Sprintf("%.2f", messageList.Amount)
					kewajiban := fmt.Sprintf("%.2f", messageList.Amount2)
					modal := fmt.Sprintf("%.2f", messageList.Amount3)
					laba := fmt.Sprintf("%.2f", LabaRugi)
					pendapatan := fmt.Sprintf("%.2f", messageList.Amount4)
					biaya := fmt.Sprintf("%.2f", messageList.Amount5)
					SendAPIPost("setPerkembanganAsetTahun" + "/" + strconv.Itoa(i) + "/" + aset + "/" + kewajiban + "/" + modal + "/" + laba + "/" + pendapatan + "/" + biaya + "/" + KodeKantor)

				}
			}

		}
		//END OF CHART PERKEMBANGAN ASET

	}

	//CHART PEMBUKAAN TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delOpenTabunganTahun" + "/" + "All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_open_tabungan_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_open_tabungan_tahun " +
			"as " +
			"SELECT t.no_rekening,pokok " +
			"from tabtrans t,tabung m,nasabah n " +
			"where t.no_rekening=m.no_rekening " +
			"and n.nasabah_id=m.nasabah_id " +
			"and floor(my_kode_trans/100)=1 " +
			"AND m.kode_produk not in ('ABP','AKP','KPE') " +
			"and t.tgl_trans=m.tgl_register " +
			"and year(m.tgl_register)= '" + s + "' ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_tabungan_tahun "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)

			SendAPIPost("setOpenTabunganTahun" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}

	//END OF CHART PEMBUKAAN TABUNGAN

	//CHART PEMBUKAAN DEPOSITO
	///DELETE DULU DATA
	SendAPIDelete("delOpenDepositoTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_open_deposito_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_open_deposito_tahun " +
			"as " +
			"SELECT t.no_rekening,pokok_trans pokok " +
			"from deptrans t,deposito m,nasabah n " +
			"where t.no_rekening=m.no_rekening " +
			"and n.nasabah_id=m.nasabah_id " +
			"and floor(my_kode_trans/100)=1 " +
			"and t.tgl_trans=m.tgl_registrasi " +
			"and year(m.tgl_registrasi)= '" + s + "' ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_open_deposito_tahun "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			SendAPIPost("setOpenDepositoTahun" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}
	//END OF CHART PEMBUKAAN DEPOSITO

	//CHART PEMBUKAAN KREDIT
	///DELETE DULU DATA
	SendAPIDelete("delOpenKreditTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_penyaluran_kredit_tahun ")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_penyaluran_kredit_tahun " +
			"as " +
			"SELECT t.no_rekening,pokok " +
			"from kretrans t,kredit m,nasabah n " +
			"where t.no_rekening=m.no_rekening " +
			"and n.nasabah_id=m.nasabah_id " +
			"and floor(my_kode_trans/100)=1 " +
			"and t.tgl_trans=m.tgl_realisasi " +
			"and year(m.tgl_realisasi)= '" + s + "' ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "SELECT coalesce(sum(pokok),0) pokok, count(no_rekening) jml from cdashboard_chart_penyaluran_kredit_tahun "
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			SendAPIPost("setOpenKreditTahun" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}

	//END OF CHART PEMBUKAAN KREDIT

	//CHART SALDO TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delSaldoTabunganTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)
		newdate := s + "-" + "12" + "-31"
		t, _ := time.Parse("2006-01-02", newdate)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			ds := strconv.Itoa(z)
			d := ""
			if z < 10 {
				d = "0" + ds
			} else {
				d = ds
			}
			z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
			ds = strconv.Itoa(z)
			dx := ""
			if z < 10 {
				dx = "0" + ds
			} else {
				dx = ds
			}
			newdate := s + "-" + d + "-" + dx
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_tabungan_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".tabung_" + cNewDate + " limit 1"
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_saldo_tabungan_tahun " +
				"as " +
				"SELECT no_rekening,saldo_akhir amount from " + DBRe + ".tabung_" + cNewDate + " " +
				"where no_rekening not in (select no_rekening from tabung where kode_produk in('ABP','AKP','KPE'))  ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_tabungan_tahun"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setSaldoTabunganTahun" + "/" + z + "/" + amount + "/" + count + "/All")

			}

		}

	}
	//END OF CHART SALDO TABUNGAN

	//CHART SALDO DEPOSITO
	///DELETE DULU DATA

	SendAPIDelete("delSaldoDepositoTahun" + "/All")
	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)
		newdate := s + "-" + "12" + "-31"
		t, _ := time.Parse("2006-01-02", newdate)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			ds := strconv.Itoa(z)
			d := ""
			if z < 10 {
				d = "0" + ds
			} else {
				d = ds
			}
			z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
			ds = strconv.Itoa(z)
			dx := ""
			if z < 10 {
				dx = "0" + ds
			} else {
				dx = ds
			}
			newdate := s + "-" + d + "-" + dx
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_deposito_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".deposito_" + cNewDate + " limit 1"
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_saldo_deposito_tahun " +
				"as " +
				"SELECT no_rekening,saldo_akhir amount from " + DBRe + ".deposito_" + cNewDate + " ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_deposito_tahun"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setSaldoDepositoTahun" + "/" + z + "/" + amount + "/" + count + "/All")

			}

		}

	}
	//END OF CHART SALDO DEPOSITO

	//CHART SALDO KREDIT

	///DELETE DULU DATA

	SendAPIDelete("delSaldoKreditTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)
		newdate := s + "-" + "12" + "-31"
		t, _ := time.Parse("2006-01-02", newdate)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			ds := strconv.Itoa(z)
			d := ""
			if z < 10 {
				d = "0" + ds
			} else {
				d = ds
			}
			z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
			ds = strconv.Itoa(z)
			dx := ""
			if z < 10 {
				dx = "0" + ds
			} else {
				dx = ds
			}
			newdate := s + "-" + d + "-" + dx
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_kredit_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where kredit.no_rekening = a.no_rekening limit 1 "
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_saldo_kredit_tahun " +
				"as " +
				"SELECT a.no_rekening,a.baki_debet amount from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
				"where a.no_rekening=kredit.no_rekening ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_kredit_tahun"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				SendAPIPost("setSaldoKreditTahun" + "/" + z + "/" + amount + "/" + count + "/All")

			}

		}

	}
	//END OF CHART SALDO KREDIT

	//CHART PERTUMBUHAN ANGSURAN

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganAngsuranTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)
		newdate := s + "-" + "12" + "-31"
		t, _ := time.Parse("2006-01-02", newdate)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			ds := strconv.Itoa(z)
			d := ""
			if z < 10 {
				d = "0" + ds
			} else {
				d = ds
			}
			z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
			ds = strconv.Itoa(z)
			dx := ""
			if z < 10 {
				dx = "0" + ds
			} else {
				dx = ds
			}
			newdate := s + "-" + d + "-" + dx
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}

		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_angsuran_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		_, err = db.Query("create view cdashboard_chart_pertumbuhan_angsuran_tahun " +
			"as " +
			"SELECT t.no_rekening, " +
			"coalesce(SUM(if(floor(my_kode_trans/100)=3,pokok,0)),0) amount," +
			"coalesce(SUM(if(floor(my_kode_trans/100)=3,bunga,0)),0) amount_bunga " +
			"FROM " +
			"kretrans t," +
			"nasabah n," +
			"kredit m " +
			"WHERE " +
			"t.no_rekening = m.no_rekening " +
			"AND n.nasabah_id = m.nasabah_id " +
			"AND tgl_trans>=date_add(date_add(LAST_DAY(last_day('" + cNewDate + "')),interval 1 DAY),interval -1 MONTH)  " +
			"AND tgl_trans<=last_day('" + cNewDate + "') group by t.no_rekening having amount<>0 or amount_bunga<>0 ")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select coalesce(sum(amount),0) amount,coalesce(sum(amount_bunga),0) amount_bunga from cdashboard_chart_pertumbuhan_angsuran_tahun"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Amount, &messageList.Amount2)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			z := strconv.Itoa(i)
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Amount2)
			SendAPIPost("setPerkembanganAngsuranTahun" + "/" + z + "/" + amount + "/" + count + "/All")

		}

	}
	//END OF CHART PERTUMBUHAN ANGSURAN

	//CHART PERKEMBANGAN NPL

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganNPLTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)
		newdate := s + "-" + "12" + "-31"
		t, _ := time.Parse("2006-01-02", newdate)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			ds := strconv.Itoa(z)
			d := ""
			if z < 10 {
				d = "0" + ds
			} else {
				d = ds
			}
			z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
			ds = strconv.Itoa(z)
			dx := ""
			if z < 10 {
				dx = "0" + ds
			} else {
				dx = ds
			}
			newdate := s + "-" + d + "-" + dx
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_perkembangan_npl_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening limit 1 "
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {

			_, err := db.Query("create view cdashboard_chart_perkembangan_npl_tahun " +
				"as " +
				"SELECT coalesce(if(sum(a.baki_debet)>0," +
				"sum(if(a.my_kolek_number>2,a.baki_debet,0))/sum(a.baki_debet)*100,0),0) as npl from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
				"where  a.no_rekening=kredit.no_rekening ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(npl,0) npl from cdashboard_chart_perkembangan_npl_tahun"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				z := strconv.Itoa(i)
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				SendAPIPost("setPerkembanganNPLTahun" + "/" + z + "/" + amount + "/All")

			}

		}

	}
	//END OF CHART PERKEMBANGAN NPL

	//CHART PERKEMBANGAN ASET,KEWAJIBAN,MODAL,LABARUGI

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganAsetTahun" + "/All")

	for i := m - 5; i <= m; i++ {
		s := strconv.Itoa(i)
		newdate := s + "-" + "12" + "-31"
		t, _ := time.Parse("2006-01-02", newdate)
		cNewDate := jodaTime.Format("YYMMdd", t)
		if m == i {
			z, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			ds := strconv.Itoa(z)
			d := ""
			if z < 10 {
				d = "0" + ds
			} else {
				d = ds
			}
			z, _ = strconv.Atoi(jodaTime.Format("dd", time.Now()))
			ds = strconv.Itoa(z)
			dx := ""
			if z < 10 {
				dx = "0" + ds
			} else {
				dx = ds
			}
			newdate := s + "-" + d + "-" + dx
			t, _ := time.Parse("2006-01-02", newdate)
			t = t.AddDate(0, 0, -1)
			cNewDate = jodaTime.Format("YYMMdd", t)
		}
		///DELETE DULU VIEW
		_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_aset_tahun ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1 "
		rowsX, err := db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_chart_pertumbuhan_aset_tahun " +
				"as " +
				"SELECT sum(if(kode_perk='1',saldo_akhir,0)) aset,sum(if(kode_perk='2',saldo_akhir,0)) kewajiban," +
				"sum(if(kode_perk='3',saldo_akhir,0)) modal,sum(if(kode_perk='4',saldo_akhir,0)) pendapatan," +
				"sum(if(kode_perk='5',saldo_akhir,0)) biaya " +
				"from " + DBRe + ".acc_" + cNewDate + " ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(aset,0) aset," +
				" coalesce(kewajiban,0) kewajiban," +
				"coalesce(modal,0) modal,coalesce(pendapatan,0) pendapatan," +
				"coalesce(biaya,0) biaya from cdashboard_chart_pertumbuhan_aset_tahun"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Amount2, &messageList.Amount3, &messageList.Amount4, &messageList.Amount5)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				LabaRugi := messageList.Amount4 - messageList.Amount5
				z := strconv.Itoa(i)
				aset := fmt.Sprintf("%.2f", messageList.Amount)
				kewajiban := fmt.Sprintf("%.2f", messageList.Amount2)
				modal := fmt.Sprintf("%.2f", messageList.Amount3)
				laba := fmt.Sprintf("%.2f", LabaRugi)
				pendapatan := fmt.Sprintf("%.2f", messageList.Amount4)
				biaya := fmt.Sprintf("%.2f", messageList.Amount5)
				SendAPIPost("setPerkembanganAsetTahun" + "/" + z + "/" + aset + "/" + kewajiban + "/" + modal + "/" + laba + "/" + pendapatan + "/" + biaya + "/All")

			}

		}

	}
	//END OF CHART PERKEMBANGAN ASET
	defer db.Close()
	functions.Logger().Info("Successfully Get Data Dashboard Chart Tahunan ")
}

func GetDataDashboardChartBulan() {
	functions.Logger().Info("Starting Scheduler Get Data Dashboard Chart Bulanan ")

	//cYear := jodaTime.Format("YYYY", time.Now())
	DBRe := os.Getenv("DB_DATABASE_RE")
	m, _ := strconv.Atoi(jodaTime.Format("YYYY", time.Now()))

	sqlStatement := "SELECT kode_kantor from app_kode_kantor order by kode_kantor "
	rowsY, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsY.Next() {
		messageList := SelectDashboardReal{}
		err = rowsY.Scan(&messageList.KodeKantor)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		KodeKantor := messageList.KodeKantor

		//CHART PEMBUKAAN TABUNGAN
		///DELETE DULU DATA
		///DELETE DULU VIEW

		m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))

		SendAPIDelete("delOpenTabunganBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))
		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			sqlStatement := "SELECT mo.month bulan,coalesce(sum( pokok ),0) pokok ,count( t.no_rekening ) jml " +
				"FROM " +
				"(SELECT 1 AS MONTH UNION SELECT 2 AS MONTH UNION SELECT 3 AS MONTH UNION SELECT 4 AS MONTH UNION SELECT 5 AS MONTH UNION SELECT 6 AS MONTH " +
				"UNION SELECT 7 AS MONTH UNION SELECT 8 AS MONTH UNION SELECT 9 AS MONTH UNION SELECT 10 AS MONTH UNION SELECT 11 AS MONTH UNION SELECT 12 AS MONTH) mo " +
				"left join tabung m on mo.month=month(m.tgl_register) AND m.kode_produk NOT IN ( 'ABP', 'AKP', 'KPE' ) " +
				"AND YEAR ( m.tgl_register )= '" + s + "' AND m.kode_kantor = '" + KodeKantor + "' " +
				"left join tabtrans t on t.no_rekening = m.no_rekening AND floor( my_kode_trans / 100 )= 1 AND t.tgl_trans = m.tgl_register " +
				"left join nasabah n  on n.nasabah_id = m.nasabah_id " +
				"group by mo.month"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Period, &messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				bulan := messageList.Period
				SendAPIPost("setOpenTabunganBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + bulan)

			}
		}
		//END OF CHART PEMBUKAAN TABUNGAN

		//CHART PEMBUKAAN DEPOSITO
		///DELETE DULU DATA
		SendAPIDelete("delOpenDepositoBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			sqlStatement := "SELECT mo.month bulan,coalesce(sum( pokok_trans ),0) pokok ,count( t.no_rekening ) jml " +
				"FROM " +
				"(SELECT 1 AS MONTH UNION SELECT 2 AS MONTH UNION SELECT 3 AS MONTH UNION SELECT 4 AS MONTH UNION SELECT 5 AS MONTH UNION SELECT 6 AS MONTH " +
				"UNION SELECT 7 AS MONTH UNION SELECT 8 AS MONTH UNION SELECT 9 AS MONTH UNION SELECT 10 AS MONTH UNION SELECT 11 AS MONTH UNION SELECT 12 AS MONTH) mo " +
				"left join deposito m on mo.month=month(m.tgl_registrasi) AND m.kode_produk NOT IN ( 'ABP', 'AKP', 'KPE' ) " +
				"AND YEAR ( m.tgl_registrasi )= '" + s + "' AND m.kode_kantor = '" + KodeKantor + "' " +
				"left join deptrans t on t.no_rekening = m.no_rekening AND floor( my_kode_trans / 100 )= 1 AND t.tgl_trans = m.tgl_registrasi " +
				"left join nasabah n  on n.nasabah_id = m.nasabah_id " +
				"group by mo.month"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Period, &messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				bulan := messageList.Period
				SendAPIPost("setOpenDepositoBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + bulan)

			}
		}
		//END OF CHART PEMBUKAAN DEPOSITO

		//CHART PEMBUKAAN KREDIT
		///DELETE DULU DATA
		SendAPIDelete("delOpenKreditBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			sqlStatement := "SELECT mo.month bulan,coalesce(sum( pokok ),0) pokok,count( t.no_rekening ) jml " +
				"FROM " +
				"(SELECT 1 AS MONTH UNION SELECT 2 AS MONTH UNION SELECT 3 AS MONTH UNION SELECT 4 AS MONTH UNION SELECT 5 AS MONTH UNION SELECT 6 AS MONTH " +
				"UNION SELECT 7 AS MONTH UNION SELECT 8 AS MONTH UNION SELECT 9 AS MONTH UNION SELECT 10 AS MONTH UNION SELECT 11 AS MONTH UNION SELECT 12 AS MONTH) mo " +
				"left join kredit m on mo.month=month(m.tgl_realisasi) " +
				"AND YEAR ( m.tgl_realisasi )= '" + s + "' AND m.kode_kantor = '" + KodeKantor + "' " +
				"left join kretrans t on t.no_rekening = m.no_rekening AND floor( my_kode_trans / 100 )= 1 AND t.tgl_trans = m.tgl_realisasi " +
				"left join nasabah n  on n.nasabah_id = m.nasabah_id " +
				"group by mo.month"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Period, &messageList.Amount, &messageList.Count)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Count)
				bulan := messageList.Period
				SendAPIPost("setOpenKreditBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + bulan)

			}
		}
		//END OF CHART PEMBUKAAN KREDIT

		//CHART SALDO TABUNGAN
		///DELETE DULU DATA
		SendAPIDelete("delSaldoTabunganBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			if m == i {
				mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
			} else {
				mm = 12
			}
			for x := 1; x <= mm; x++ {
				dd := "31"
				if x == 1 {
					dd = "31"
				} else if x == 3 {
					dd = "31"
				} else if x == 4 {
					dd = "30"
				} else if x == 5 {
					dd = "31"
				} else if x == 6 {
					dd = "30"
				} else if x == 7 {
					dd = "31"
				} else if x == 8 {
					dd = "31"
				} else if x == 9 {
					dd = "30"
				} else if x == 10 {
					dd = "31"
				} else if x == 11 {
					dd = "30"
				} else if x == 12 {
					dd = "31"
				}
				xs := ""
				if x < 10 {
					xs = "0" + strconv.Itoa(x)
				} else {
					xs = strconv.Itoa(x)
				}
				newdate := s + "-" + xs + "-" + dd
				if x == 2 {
					newdate = s + "-" + "03" + "-" + "01"
				}
				t, _ := time.Parse("2006-01-02", newdate)
				if x == 2 {
					t = t.AddDate(0, 0, -1)
				}
				cNewDate := jodaTime.Format("YYMMdd", t)
				///DELETE DULU VIEW
				_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_tabungan_bulan ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatementX := "select 1 from " + DBRe + ".tabung_" + cNewDate + " a,tabung where tabung.no_Rekening=a.no_rekening and tabung.kode_kantor='" + KodeKantor + "' limit 1"
				rowsX, err := db.Query(sqlStatementX)
				if err != nil {
					SendAPIPost("setSaldoTabunganBulan" + "/" + s + "/" + "0.00" + "/" + "0" + "/" + KodeKantor + "/" + strconv.Itoa(x))
				}
				if err == nil {
					for rowsX.Next() {
						_, err = db.Query("create view cdashboard_chart_saldo_tabungan_bulan " +
							"as " +
							"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".tabung_" + cNewDate + " a,tabung " +
							"where tabung.no_rekening=a.no_rekening and tabung.kode_kantor='" + KodeKantor + "' " +
							"AND tabung.kode_produk not in ('ABP','AKP','KPE') ")

						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_tabungan_bulan"
						rows, err := db.Query(sqlStatement)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						for rows.Next() {
							messageList := SelectDashboardChart{}
							err = rows.Scan(&messageList.Amount, &messageList.Count)
							if err != nil {
								functions.Logger().Error(err.Error())
								return
							}

							amount := fmt.Sprintf("%.2f", messageList.Amount)
							count := fmt.Sprintf("%.2f", messageList.Count)
							SendAPIPost("setSaldoTabunganBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + strconv.Itoa(x))
						}
					}
				}
			}
		}
		//END OF CHART SALDO TABUNGAN

		//CHART SALDO DEPOSITO
		///DELETE DULU DATA

		SendAPIDelete("delSaldoDepositoBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			if m == i {
				mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
			} else {
				mm = 12
			}
			for x := 1; x <= mm; x++ {
				dd := "31"
				if x == 1 {
					dd = "31"
				} else if x == 3 {
					dd = "31"
				} else if x == 4 {
					dd = "30"
				} else if x == 5 {
					dd = "31"
				} else if x == 6 {
					dd = "30"
				} else if x == 7 {
					dd = "31"
				} else if x == 8 {
					dd = "31"
				} else if x == 9 {
					dd = "30"
				} else if x == 10 {
					dd = "31"
				} else if x == 11 {
					dd = "30"
				} else if x == 12 {
					dd = "31"
				}
				xs := ""
				if x < 10 {
					xs = "0" + strconv.Itoa(x)
				} else {
					xs = strconv.Itoa(x)
				}
				newdate := s + "-" + xs + "-" + dd
				if x == 2 {
					newdate = s + "-" + "03" + "-" + "01"
				}
				t, _ := time.Parse("2006-01-02", newdate)
				if x == 2 {
					t = t.AddDate(0, 0, -1)
				}
				cNewDate := jodaTime.Format("YYMMdd", t)
				///DELETE DULU VIEW
				_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_deposito_bulan ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatementX := "select 1 from " + DBRe + ".deposito_" + cNewDate + " a,deposito" +
					" where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' limit 1"
				rowsX, err := db.Query(sqlStatementX)
				if err != nil {
					SendAPIPost("setSaldoDepositoBulan" + "/" + s + "/" + "0.00" + "/" + "0" + "/" + KodeKantor + "/" + strconv.Itoa(x))
				}
				if err == nil {
					for rowsX.Next() {
						_, err := db.Query("create view cdashboard_chart_saldo_deposito_bulan " +
							"as " +
							"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".deposito_" + cNewDate + " a,deposito " +
							"where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' ")

						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_deposito_bulan"
						rows, err := db.Query(sqlStatement)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						for rows.Next() {
							messageList := SelectDashboardChart{}
							err = rows.Scan(&messageList.Amount, &messageList.Count)
							if err != nil {
								functions.Logger().Error(err.Error())
								return
							}

							amount := fmt.Sprintf("%.2f", messageList.Amount)
							count := fmt.Sprintf("%.2f", messageList.Count)
							SendAPIPost("setSaldoDepositoBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + strconv.Itoa(x))
						}
					}
				}
			}
		}

		SendAPIDelete("delSaldoKreditBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			if m == i {
				mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
			} else {
				mm = 12
			}
			for x := 1; x <= mm; x++ {
				dd := "31"
				if x == 1 {
					dd = "31"
				} else if x == 3 {
					dd = "31"
				} else if x == 4 {
					dd = "30"
				} else if x == 5 {
					dd = "31"
				} else if x == 6 {
					dd = "30"
				} else if x == 7 {
					dd = "31"
				} else if x == 8 {
					dd = "31"
				} else if x == 9 {
					dd = "30"
				} else if x == 10 {
					dd = "31"
				} else if x == 11 {
					dd = "30"
				} else if x == 12 {
					dd = "31"
				}
				xs := ""
				if x < 10 {
					xs = "0" + strconv.Itoa(x)
				} else {
					xs = strconv.Itoa(x)
				}
				newdate := s + "-" + xs + "-" + dd
				if x == 2 {
					newdate = s + "-" + "03" + "-" + "01"
				}
				t, _ := time.Parse("2006-01-02", newdate)
				if x == 2 {
					t = t.AddDate(0, 0, -1)
				}
				cNewDate := jodaTime.Format("YYMMdd", t)
				///DELETE DULU VIEW
				_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_kredit_bulan ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a, kredit " +
					"where kredit.no_rekening=a.no_rekening and kredit.kode_kantor='" + KodeKantor + "'  limit 1"
				rowsX, err := db.Query(sqlStatementX)
				if err != nil {
					SendAPIPost("setSaldoKreditBulan" + "/" + s + "/" + "0.00" + "/" + "0" + "/" + KodeKantor + "/" + strconv.Itoa(x))
				}
				if err == nil {
					for rowsX.Next() {
						_, err := db.Query("create view cdashboard_chart_saldo_kredit_bulan " +
							"as " +
							"SELECT a.no_rekening,a.baki_debet amount from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening " +
							"and kredit.kode_kantor='" + KodeKantor + "'  ")

						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_kredit_bulan"
						rows, err := db.Query(sqlStatement)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						for rows.Next() {
							messageList := SelectDashboardChart{}
							err = rows.Scan(&messageList.Amount, &messageList.Count)
							if err != nil {
								functions.Logger().Error(err.Error())
								return
							}

							amount := fmt.Sprintf("%.2f", messageList.Amount)
							count := fmt.Sprintf("%.2f", messageList.Count)
							SendAPIPost("setSaldoKreditBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + strconv.Itoa(x))
						}
					}
				}
			}
		}

		SendAPIDelete("delPerkembanganAngsuranBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			if m == i {
				mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
			} else {
				mm = 12
			}
			for x := 1; x <= mm; x++ {
				dd := "31"
				if x == 1 {
					dd = "31"
				} else if x == 3 {
					dd = "31"
				} else if x == 4 {
					dd = "30"
				} else if x == 5 {
					dd = "31"
				} else if x == 6 {
					dd = "30"
				} else if x == 7 {
					dd = "31"
				} else if x == 8 {
					dd = "31"
				} else if x == 9 {
					dd = "30"
				} else if x == 10 {
					dd = "31"
				} else if x == 11 {
					dd = "30"
				} else if x == 12 {
					dd = "31"
				}
				xs := ""
				if x < 10 {
					xs = "0" + strconv.Itoa(x)
				} else {
					xs = strconv.Itoa(x)
				}
				newdate := s + "-" + xs + "-" + dd
				if x == 2 {
					newdate = s + "-" + "03" + "-" + "01"
				}
				t, _ := time.Parse("2006-01-02", newdate)
				if x == 2 {
					t = t.AddDate(0, 0, -1)
				}
				cNewDate := jodaTime.Format("YYMMdd", t)
				///DELETE DULU VIEW
				_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_angsuran_bulan ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				_, err = db.Query("create view cdashboard_chart_pertumbuhan_angsuran_bulan " +
					"as " +
					"SELECT t.no_rekening, " +
					"coalesce(SUM(if(floor(my_kode_trans/100)=3,pokok,0)),0) amount," +
					"coalesce(SUM(if(floor(my_kode_trans/100)=3,bunga,0)),0) amount_bunga " +
					"FROM " +
					"kretrans t," +
					"nasabah n," +
					"kredit m " +
					"WHERE " +
					"t.no_rekening = m.no_rekening " +
					"AND n.nasabah_id = m.nasabah_id " +
					"AND tgl_trans>=date_add(date_add(LAST_DAY(last_day('" + cNewDate + "')),interval 1 DAY),interval -1 MONTH)  " +
					"AND tgl_trans<=last_day('" + cNewDate + "') and m.kode_kantor='" + KodeKantor + "' group by t.no_rekening having amount<>0 or amount_bunga<>0 ")

				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatement := "select coalesce(sum(amount),0) amount,coalesce(sum(amount_bunga),0) amount_bunga from cdashboard_chart_pertumbuhan_angsuran_bulan"
				rows, err := db.Query(sqlStatement)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				for rows.Next() {
					messageList := SelectDashboardChart{}
					err = rows.Scan(&messageList.Amount, &messageList.Amount2)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					amount := fmt.Sprintf("%.2f", messageList.Amount)
					count := fmt.Sprintf("%.2f", messageList.Amount2)
					SendAPIPost("setPerkembanganAngsuranBulan" + "/" + s + "/" + amount + "/" + count + "/" + KodeKantor + "/" + strconv.Itoa(x))
				}
			}
		}
		//END OF CHART SALDO KREDIT

		//CHART PERKEMBANGAN NPL

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganNPLBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			if m == i {
				mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
			} else {
				mm = 12
			}
			for x := 1; x <= mm; x++ {
				dd := "31"
				if x == 1 {
					dd = "31"
				} else if x == 3 {
					dd = "31"
				} else if x == 4 {
					dd = "30"
				} else if x == 5 {
					dd = "31"
				} else if x == 6 {
					dd = "30"
				} else if x == 7 {
					dd = "31"
				} else if x == 8 {
					dd = "31"
				} else if x == 9 {
					dd = "30"
				} else if x == 10 {
					dd = "31"
				} else if x == 11 {
					dd = "30"
				} else if x == 12 {
					dd = "31"
				}
				xs := ""
				if x < 10 {
					xs = "0" + strconv.Itoa(x)
				} else {
					xs = strconv.Itoa(x)
				}
				newdate := s + "-" + xs + "-" + dd
				if x == 2 {
					newdate = s + "-" + "03" + "-" + "01"
				}
				t, _ := time.Parse("2006-01-02", newdate)
				if x == 2 {
					t = t.AddDate(0, 0, -1)
				}
				cNewDate := jodaTime.Format("YYMMdd", t)
				///DELETE DULU VIEW
				_, err := db.Query("drop view IF EXISTS cdashboard_chart_perkembangan_npl_tahun ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "' limit 1"
				rowsX, err := db.Query(sqlStatementX)
				if err != nil {
					SendAPIPost("setPerkembanganNPLBulan" + "/" + s + "/" + "0.00" + "/" + KodeKantor + "/" + strconv.Itoa(x))
				}
				if err == nil {
					for rowsX.Next() {

						_, err := db.Query("create view cdashboard_chart_perkembangan_npl_tahun " +
							"as " +
							"SELECT coalesce(if(sum(a.baki_debet)>0," +
							"sum(if(a.my_kolek_number>2,a.baki_debet,0))/sum(a.baki_debet)*100,0),0) as npl from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
							"where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "'")

						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						sqlStatement := "select coalesce(npl,0) npl from cdashboard_chart_perkembangan_npl_tahun"
						rows, err := db.Query(sqlStatement)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						for rows.Next() {
							messageList := SelectDashboardChart{}
							err = rows.Scan(&messageList.Amount)
							if err != nil {
								functions.Logger().Error(err.Error())
								return
							}
							amount := fmt.Sprintf("%.2f", messageList.Amount)
							SendAPIPost("setPerkembanganNPLBulan" + "/" + s + "/" + amount + "/" + KodeKantor + "/" + strconv.Itoa(x))

						}
					}
				}
			}
		}
		//END OF CHART PERKEMBANGAN NPL

		//CHART PERKEMBANGAN ASET,KEWAJIBAN,MODAL,LABARUGI

		///DELETE DULU DATA
		SendAPIDelete("delPerkembanganAsetBulan" + "/" + KodeKantor + "/" + strconv.Itoa(m))

		for i := m; i <= m; i++ {
			s := strconv.Itoa(i)

			mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
			if m == i {
				mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
			} else {
				mm = 12
			}
			for x := 1; x <= mm; x++ {
				dd := "31"
				if x == 1 {
					dd = "31"
				} else if x == 3 {
					dd = "31"
				} else if x == 4 {
					dd = "30"
				} else if x == 5 {
					dd = "31"
				} else if x == 6 {
					dd = "30"
				} else if x == 7 {
					dd = "31"
				} else if x == 8 {
					dd = "31"
				} else if x == 9 {
					dd = "30"
				} else if x == 10 {
					dd = "31"
				} else if x == 11 {
					dd = "30"
				} else if x == 12 {
					dd = "31"
				}
				xs := ""
				if x < 10 {
					xs = "0" + strconv.Itoa(x)
				} else {
					xs = strconv.Itoa(x)
				}
				newdate := s + "-" + xs + "-" + dd
				if x == 2 {
					newdate = s + "-" + "03" + "-" + "01"
				}
				t, _ := time.Parse("2006-01-02", newdate)
				if x == 2 {
					t = t.AddDate(0, 0, -1)
				}
				cNewDate := jodaTime.Format("YYMMdd", t)
				///DELETE DULU VIEW
				_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_aset_tahun ")
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
				rowsX, err := db.Query(sqlStatementX)
				if err != nil {
					SendAPIPost("setPerkembanganAsetBulan" + "/" + s + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + KodeKantor + "/" + strconv.Itoa(x))
				}
				if err == nil {
					for rowsX.Next() {
						_, err := db.Query("create view cdashboard_chart_pertumbuhan_aset_tahun " +
							"as " +
							"SELECT sum(if(kode_perk='1',saldo_akhir,0)) aset,sum(if(kode_perk='2',saldo_akhir,0)) kewajiban," +
							"sum(if(kode_perk='3',saldo_akhir,0)) modal,sum(if(kode_perk='4',saldo_akhir,0)) pendapatan," +
							"sum(if(kode_perk='5',saldo_akhir,0)) biaya " +
							"from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' ")
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						sqlStatement := "select coalesce(aset,0) aset," +
							" coalesce(kewajiban,0) kewajiban," +
							"coalesce(modal,0) modal,coalesce(pendapatan,0) pendapatan," +
							"coalesce(biaya,0) biaya from cdashboard_chart_pertumbuhan_aset_tahun"
						rows, err := db.Query(sqlStatement)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						for rows.Next() {
							messageList := SelectDashboardChart{}
							err = rows.Scan(&messageList.Amount, &messageList.Amount2, &messageList.Amount3, &messageList.Amount4, &messageList.Amount5)
							if err != nil {
								functions.Logger().Error(err.Error())
								return
							}

							LabaRugi := messageList.Amount4 - messageList.Amount5
							aset := fmt.Sprintf("%.2f", messageList.Amount)
							kewajiban := fmt.Sprintf("%.2f", messageList.Amount2)
							modal := fmt.Sprintf("%.2f", messageList.Amount3)
							laba := fmt.Sprintf("%.2f", LabaRugi)
							pendapatan := fmt.Sprintf("%.2f", messageList.Amount4)
							biaya := fmt.Sprintf("%.2f", messageList.Amount5)
							SendAPIPost("setPerkembanganAsetBulan" + "/" + s + "/" + aset + "/" + kewajiban + "/" + modal + "/" + laba + "/" + pendapatan + "/" + biaya + "/" + KodeKantor + "/" + strconv.Itoa(x))

						}
					}
				}
			}
		}
		//END OF CHART PERKEMBANGAN ASET
	}

	//CHART PEMBUKAAN TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delOpenTabunganBulan" + "/" + "All" + "/" + strconv.Itoa(m))
	///DELETE DULU VIEW

	m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		sqlStatement := "SELECT mo.month bulan,coalesce(sum( pokok ),0) pokok ,count( t.no_rekening ) jml " +
			"FROM " +
			"(SELECT 1 AS MONTH UNION SELECT 2 AS MONTH UNION SELECT 3 AS MONTH UNION SELECT 4 AS MONTH UNION SELECT 5 AS MONTH UNION SELECT 6 AS MONTH " +
			"UNION SELECT 7 AS MONTH UNION SELECT 8 AS MONTH UNION SELECT 9 AS MONTH UNION SELECT 10 AS MONTH UNION SELECT 11 AS MONTH UNION SELECT 12 AS MONTH) mo " +
			"left join tabung m on mo.month=month(m.tgl_register) AND m.kode_produk NOT IN ( 'ABP', 'AKP', 'KPE' ) " +
			"AND YEAR ( m.tgl_register )= '" + s + "' " +
			"left join tabtrans t on t.no_rekening = m.no_rekening AND floor( my_kode_trans / 100 )= 1 AND t.tgl_trans = m.tgl_register " +
			"left join nasabah n  on n.nasabah_id = m.nasabah_id " +
			"group by mo.month"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Period, &messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			bulan := messageList.Period
			SendAPIPost("setOpenTabunganBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + bulan)

		}
	}
	//END OF CHART PEMBUKAAN TABUNGAN

	//CHART PEMBUKAAN DEPOSITO
	///DELETE DULU DATA
	SendAPIDelete("delOpenDepositoBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		sqlStatement := "SELECT mo.month bulan,coalesce(sum( pokok_trans ),0) pokok ,count( t.no_rekening ) jml " +
			"FROM " +
			"(SELECT 1 AS MONTH UNION SELECT 2 AS MONTH UNION SELECT 3 AS MONTH UNION SELECT 4 AS MONTH UNION SELECT 5 AS MONTH UNION SELECT 6 AS MONTH " +
			"UNION SELECT 7 AS MONTH UNION SELECT 8 AS MONTH UNION SELECT 9 AS MONTH UNION SELECT 10 AS MONTH UNION SELECT 11 AS MONTH UNION SELECT 12 AS MONTH) mo " +
			"left join deposito m on mo.month=month(m.tgl_registrasi) AND m.kode_produk NOT IN ( 'ABP', 'AKP', 'KPE' ) " +
			"AND YEAR ( m.tgl_registrasi )= '" + s + "'  " +
			"left join deptrans t on t.no_rekening = m.no_rekening AND floor( my_kode_trans / 100 )= 1 AND t.tgl_trans = m.tgl_registrasi " +
			"left join nasabah n  on n.nasabah_id = m.nasabah_id " +
			"group by mo.month"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Period, &messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			bulan := messageList.Period
			SendAPIPost("setOpenDepositoBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + bulan)

		}
	}
	//END OF CHART PEMBUKAAN DEPOSITO

	//CHART PEMBUKAAN KREDIT
	///DELETE DULU DATA
	SendAPIDelete("delOpenKreditBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	m, _ = strconv.Atoi(jodaTime.Format("YYYY", time.Now()))
	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		sqlStatement := "SELECT mo.month bulan,coalesce(sum( pokok ),0) pokok,count( t.no_rekening ) jml " +
			"FROM " +
			"(SELECT 1 AS MONTH UNION SELECT 2 AS MONTH UNION SELECT 3 AS MONTH UNION SELECT 4 AS MONTH UNION SELECT 5 AS MONTH UNION SELECT 6 AS MONTH " +
			"UNION SELECT 7 AS MONTH UNION SELECT 8 AS MONTH UNION SELECT 9 AS MONTH UNION SELECT 10 AS MONTH UNION SELECT 11 AS MONTH UNION SELECT 12 AS MONTH) mo " +
			"left join kredit m on mo.month=month(m.tgl_realisasi) " +
			"AND YEAR ( m.tgl_realisasi )= '" + s + "' " +
			"left join kretrans t on t.no_rekening = m.no_rekening AND floor( my_kode_trans / 100 )= 1 AND t.tgl_trans = m.tgl_realisasi " +
			"left join nasabah n  on n.nasabah_id = m.nasabah_id " +
			"group by mo.month"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardChart{}
			err = rows.Scan(&messageList.Period, &messageList.Amount, &messageList.Count)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			amount := fmt.Sprintf("%.2f", messageList.Amount)
			count := fmt.Sprintf("%.2f", messageList.Count)
			bulan := messageList.Period
			SendAPIPost("setOpenKreditBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + bulan)

		}
	}
	//END OF CHART PEMBUKAAN KREDIT
	//CHART SALDO TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delSaldoTabunganBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		if m == i {
			mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		} else {
			mm = 12
		}
		for x := 1; x <= mm; x++ {
			dd := "31"
			if x == 1 {
				dd = "31"
			} else if x == 3 {
				dd = "31"
			} else if x == 4 {
				dd = "30"
			} else if x == 5 {
				dd = "31"
			} else if x == 6 {
				dd = "30"
			} else if x == 7 {
				dd = "31"
			} else if x == 8 {
				dd = "31"
			} else if x == 9 {
				dd = "30"
			} else if x == 10 {
				dd = "31"
			} else if x == 11 {
				dd = "30"
			} else if x == 12 {
				dd = "31"
			}
			xs := ""
			if x < 10 {
				xs = "0" + strconv.Itoa(x)
			} else {
				xs = strconv.Itoa(x)
			}
			newdate := s + "-" + xs + "-" + dd
			if x == 2 {
				newdate = s + "-" + "03" + "-" + "01"
			}
			t, _ := time.Parse("2006-01-02", newdate)
			if x == 2 {
				t = t.AddDate(0, 0, -1)
			}
			cNewDate := jodaTime.Format("YYMMdd", t)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_tabungan_bulan ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".tabung_" + cNewDate + " a,tabung where tabung.no_Rekening=a.no_rekening  limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				SendAPIPost("setSaldoTabunganBulan" + "/" + s + "/" + "0.00" + "/" + "0" + "/" + "All" + "/" + strconv.Itoa(x))
			}
			if err == nil {
				for rowsX.Next() {
					_, err = db.Query("create view cdashboard_chart_saldo_tabungan_bulan " +
						"as " +
						"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".tabung_" + cNewDate + " a,tabung " +
						"where tabung.no_rekening=a.no_rekening  " +
						"AND tabung.kode_produk not in ('ABP','AKP','KPE') ")

					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_tabungan_bulan"
					rows, err := db.Query(sqlStatement)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					for rows.Next() {
						messageList := SelectDashboardChart{}
						err = rows.Scan(&messageList.Amount, &messageList.Count)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						amount := fmt.Sprintf("%.2f", messageList.Amount)
						count := fmt.Sprintf("%.2f", messageList.Count)
						SendAPIPost("setSaldoTabunganBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + strconv.Itoa(x))
					}
				}
			}
		}
	}
	//END OF CHART SALDO TABUNGAN

	//CHART SALDO DEPOSITO
	///DELETE DULU DATA

	SendAPIDelete("delSaldoDepositoBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		if m == i {
			mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		} else {
			mm = 12
		}
		for x := 1; x <= mm; x++ {
			dd := "31"
			if x == 1 {
				dd = "31"
			} else if x == 3 {
				dd = "31"
			} else if x == 4 {
				dd = "30"
			} else if x == 5 {
				dd = "31"
			} else if x == 6 {
				dd = "30"
			} else if x == 7 {
				dd = "31"
			} else if x == 8 {
				dd = "31"
			} else if x == 9 {
				dd = "30"
			} else if x == 10 {
				dd = "31"
			} else if x == 11 {
				dd = "30"
			} else if x == 12 {
				dd = "31"
			}
			xs := ""
			if x < 10 {
				xs = "0" + strconv.Itoa(x)
			} else {
				xs = strconv.Itoa(x)
			}
			newdate := s + "-" + xs + "-" + dd
			if x == 2 {
				newdate = s + "-" + "03" + "-" + "01"
			}
			t, _ := time.Parse("2006-01-02", newdate)
			if x == 2 {
				t = t.AddDate(0, 0, -1)
			}
			cNewDate := jodaTime.Format("YYMMdd", t)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_deposito_bulan ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".deposito_" + cNewDate + " a,deposito" +
				" where a.no_rekening=deposito.no_rekening limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				SendAPIPost("setSaldoDepositoBulan" + "/" + s + "/" + "0.00" + "/" + "0" + "/" + "All" + "/" + strconv.Itoa(x))
			}
			if err == nil {
				for rowsX.Next() {
					_, err := db.Query("create view cdashboard_chart_saldo_deposito_bulan " +
						"as " +
						"SELECT a.no_rekening,a.saldo_akhir amount from " + DBRe + ".deposito_" + cNewDate + " a,deposito " +
						"where a.no_rekening=deposito.no_rekening  ")

					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_deposito_bulan"
					rows, err := db.Query(sqlStatement)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					for rows.Next() {
						messageList := SelectDashboardChart{}
						err = rows.Scan(&messageList.Amount, &messageList.Count)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						amount := fmt.Sprintf("%.2f", messageList.Amount)
						count := fmt.Sprintf("%.2f", messageList.Count)
						SendAPIPost("setSaldoDepositoBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + strconv.Itoa(x))
					}
				}
			}
		}
	}

	SendAPIDelete("delSaldoKreditBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		if m == i {
			mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		} else {
			mm = 12
		}
		for x := 1; x <= mm; x++ {
			dd := "31"
			if x == 1 {
				dd = "31"
			} else if x == 3 {
				dd = "31"
			} else if x == 4 {
				dd = "30"
			} else if x == 5 {
				dd = "31"
			} else if x == 6 {
				dd = "30"
			} else if x == 7 {
				dd = "31"
			} else if x == 8 {
				dd = "31"
			} else if x == 9 {
				dd = "30"
			} else if x == 10 {
				dd = "31"
			} else if x == 11 {
				dd = "30"
			} else if x == 12 {
				dd = "31"
			}
			xs := ""
			if x < 10 {
				xs = "0" + strconv.Itoa(x)
			} else {
				xs = strconv.Itoa(x)
			}
			newdate := s + "-" + xs + "-" + dd
			if x == 2 {
				newdate = s + "-" + "03" + "-" + "01"
			}
			t, _ := time.Parse("2006-01-02", newdate)
			if x == 2 {
				t = t.AddDate(0, 0, -1)
			}
			cNewDate := jodaTime.Format("YYMMdd", t)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_saldo_kredit_bulan ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a, kredit " +
				"where kredit.no_rekening=a.no_rekening  limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				SendAPIPost("setSaldoKreditBulan" + "/" + s + "/" + "0.00" + "/" + "0" + "/" + "All" + "/" + strconv.Itoa(x))
			}
			if err == nil {
				for rowsX.Next() {
					_, err := db.Query("create view cdashboard_chart_saldo_kredit_bulan " +
						"as " +
						"SELECT a.no_rekening,a.baki_debet amount from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening " +
						" ")

					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					sqlStatement := "select coalesce(sum(amount),0) amount,count(no_rekening) count from cdashboard_chart_saldo_kredit_bulan"
					rows, err := db.Query(sqlStatement)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					for rows.Next() {
						messageList := SelectDashboardChart{}
						err = rows.Scan(&messageList.Amount, &messageList.Count)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						amount := fmt.Sprintf("%.2f", messageList.Amount)
						count := fmt.Sprintf("%.2f", messageList.Count)
						SendAPIPost("setSaldoKreditBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + strconv.Itoa(x))
					}
				}
			}
		}
	}
	SendAPIDelete("delPerkembanganAngsuranBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		if m == i {
			mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		} else {
			mm = 12
		}
		for x := 1; x <= mm; x++ {
			dd := "31"
			if x == 1 {
				dd = "31"
			} else if x == 3 {
				dd = "31"
			} else if x == 4 {
				dd = "30"
			} else if x == 5 {
				dd = "31"
			} else if x == 6 {
				dd = "30"
			} else if x == 7 {
				dd = "31"
			} else if x == 8 {
				dd = "31"
			} else if x == 9 {
				dd = "30"
			} else if x == 10 {
				dd = "31"
			} else if x == 11 {
				dd = "30"
			} else if x == 12 {
				dd = "31"
			}
			xs := ""
			if x < 10 {
				xs = "0" + strconv.Itoa(x)
			} else {
				xs = strconv.Itoa(x)
			}
			newdate := s + "-" + xs + "-" + dd
			if x == 2 {
				newdate = s + "-" + "03" + "-" + "01"
			}
			t, _ := time.Parse("2006-01-02", newdate)
			if x == 2 {
				t = t.AddDate(0, 0, -1)
			}
			cNewDate := jodaTime.Format("YYMMdd", t)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_angsuran_bulan ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			_, err = db.Query("create view cdashboard_chart_pertumbuhan_angsuran_bulan " +
				"as " +
				"SELECT t.no_rekening, " +
				"coalesce(SUM(if(floor(my_kode_trans/100)=3,pokok,0)),0) amount," +
				"coalesce(SUM(if(floor(my_kode_trans/100)=3,bunga,0)),0) amount_bunga " +
				"FROM " +
				"kretrans t," +
				"nasabah n," +
				"kredit m " +
				"WHERE " +
				"t.no_rekening = m.no_rekening " +
				"AND n.nasabah_id = m.nasabah_id " +
				"AND tgl_trans>=date_add(date_add(LAST_DAY(last_day('" + cNewDate + "')),interval 1 DAY),interval -1 MONTH)  " +
				"AND tgl_trans<=last_day('" + cNewDate + "')  group by t.no_rekening having amount<>0 or amount_bunga<>0 ")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select coalesce(sum(amount),0) amount,coalesce(sum(amount_bunga),0) amount_bunga from cdashboard_chart_pertumbuhan_angsuran_bulan"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardChart{}
				err = rows.Scan(&messageList.Amount, &messageList.Amount2)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				amount := fmt.Sprintf("%.2f", messageList.Amount)
				count := fmt.Sprintf("%.2f", messageList.Amount2)
				SendAPIPost("setPerkembanganAngsuranBulan" + "/" + s + "/" + amount + "/" + count + "/" + "All" + "/" + strconv.Itoa(x))
			}
		}
	}
	//END OF CHART SALDO KREDIT

	//CHART PERKEMBANGAN NPL

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganNPLBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		if m == i {
			mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		} else {
			mm = 12
		}
		for x := 1; x <= mm; x++ {
			dd := "31"
			if x == 1 {
				dd = "31"
			} else if x == 3 {
				dd = "31"
			} else if x == 4 {
				dd = "30"
			} else if x == 5 {
				dd = "31"
			} else if x == 6 {
				dd = "30"
			} else if x == 7 {
				dd = "31"
			} else if x == 8 {
				dd = "31"
			} else if x == 9 {
				dd = "30"
			} else if x == 10 {
				dd = "31"
			} else if x == 11 {
				dd = "30"
			} else if x == 12 {
				dd = "31"
			}
			xs := ""
			if x < 10 {
				xs = "0" + strconv.Itoa(x)
			} else {
				xs = strconv.Itoa(x)
			}
			newdate := s + "-" + xs + "-" + dd
			if x == 2 {
				newdate = s + "-" + "03" + "-" + "01"
			}
			t, _ := time.Parse("2006-01-02", newdate)
			if x == 2 {
				t = t.AddDate(0, 0, -1)
			}
			cNewDate := jodaTime.Format("YYMMdd", t)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_perkembangan_npl_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening  limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				SendAPIPost("setPerkembanganNPLBulan" + "/" + s + "/" + "0.00" + "/" + "All" + "/" + strconv.Itoa(x))
			}
			if err == nil {
				for rowsX.Next() {

					_, err := db.Query("create view cdashboard_chart_perkembangan_npl_tahun " +
						"as " +
						"SELECT coalesce(if(sum(a.baki_debet)>0," +
						"sum(if(a.my_kolek_number>2,a.baki_debet,0))/sum(a.baki_debet)*100,0),0) as npl from " + DBRe + ".kredit_" + cNewDate + " a,kredit " +
						"where a.no_rekening=kredit.no_rekening ")

					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					sqlStatement := "select coalesce(npl,0) npl from cdashboard_chart_perkembangan_npl_tahun"
					rows, err := db.Query(sqlStatement)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					for rows.Next() {
						messageList := SelectDashboardChart{}
						err = rows.Scan(&messageList.Amount)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}
						amount := fmt.Sprintf("%.2f", messageList.Amount)
						SendAPIPost("setPerkembanganNPLBulan" + "/" + s + "/" + amount + "/" + "All" + "/" + strconv.Itoa(x))

					}
				}
			}
		}
	}
	//END OF CHART PERKEMBANGAN NPL

	//CHART PERKEMBANGAN ASET,KEWAJIBAN,MODAL,LABARUGI

	///DELETE DULU DATA
	SendAPIDelete("delPerkembanganAsetBulan" + "/" + "All" + "/" + strconv.Itoa(m))

	for i := m; i <= m; i++ {
		s := strconv.Itoa(i)

		mm, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
		if m == i {
			mm, _ = strconv.Atoi(jodaTime.Format("MM", time.Now()))
		} else {
			mm = 12
		}
		for x := 1; x <= mm; x++ {
			dd := "31"
			if x == 1 {
				dd = "31"
			} else if x == 3 {
				dd = "31"
			} else if x == 4 {
				dd = "30"
			} else if x == 5 {
				dd = "31"
			} else if x == 6 {
				dd = "30"
			} else if x == 7 {
				dd = "31"
			} else if x == 8 {
				dd = "31"
			} else if x == 9 {
				dd = "30"
			} else if x == 10 {
				dd = "31"
			} else if x == 11 {
				dd = "30"
			} else if x == 12 {
				dd = "31"
			}
			xs := ""
			if x < 10 {
				xs = "0" + strconv.Itoa(x)
			} else {
				xs = strconv.Itoa(x)
			}
			newdate := s + "-" + xs + "-" + dd
			if x == 2 {
				newdate = s + "-" + "03" + "-" + "01"
			}
			t, _ := time.Parse("2006-01-02", newdate)
			if x == 2 {
				t = t.AddDate(0, 0, -1)
			}
			cNewDate := jodaTime.Format("YYMMdd", t)
			///DELETE DULU VIEW
			_, err := db.Query("drop view IF EXISTS cdashboard_chart_pertumbuhan_aset_tahun ")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + "  limit 1"
			rowsX, err := db.Query(sqlStatementX)
			if err != nil {
				SendAPIPost("setPerkembanganAsetBulan" + "/" + s + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "0.00" + "/" + "All" + "/" + strconv.Itoa(x))
			}
			if err == nil {
				for rowsX.Next() {
					_, err := db.Query("create view cdashboard_chart_pertumbuhan_aset_tahun " +
						"as " +
						"SELECT sum(if(kode_perk='1',saldo_akhir,0)) aset,sum(if(kode_perk='2',saldo_akhir,0)) kewajiban," +
						"sum(if(kode_perk='3',saldo_akhir,0)) modal,sum(if(kode_perk='4',saldo_akhir,0)) pendapatan," +
						"sum(if(kode_perk='5',saldo_akhir,0)) biaya " +
						"from " + DBRe + ".acc_" + cNewDate + "  ")
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					sqlStatement := "select coalesce(aset,0) aset," +
						" coalesce(kewajiban,0) kewajiban," +
						"coalesce(modal,0) modal,coalesce(pendapatan,0) pendapatan," +
						"coalesce(biaya,0) biaya from cdashboard_chart_pertumbuhan_aset_tahun"
					rows, err := db.Query(sqlStatement)
					if err != nil {
						functions.Logger().Error(err.Error())
						return
					}

					for rows.Next() {
						messageList := SelectDashboardChart{}
						err = rows.Scan(&messageList.Amount, &messageList.Amount2, &messageList.Amount3, &messageList.Amount4, &messageList.Amount5)
						if err != nil {
							functions.Logger().Error(err.Error())
							return
						}

						LabaRugi := messageList.Amount4 - messageList.Amount5
						aset := fmt.Sprintf("%.2f", messageList.Amount)
						kewajiban := fmt.Sprintf("%.2f", messageList.Amount2)
						modal := fmt.Sprintf("%.2f", messageList.Amount3)
						laba := fmt.Sprintf("%.2f", LabaRugi)
						pendapatan := fmt.Sprintf("%.2f", messageList.Amount4)
						biaya := fmt.Sprintf("%.2f", messageList.Amount5)
						SendAPIPost("setPerkembanganAsetBulan" + "/" + s + "/" + aset + "/" + kewajiban + "/" + modal + "/" + laba + "/" + pendapatan + "/" + biaya + "/" + "All" + "/" + strconv.Itoa(x))

					}
				}
			}
		}

	}
	defer db.Close()
	functions.Logger().Info("Successfully Get Data Dashboard Chart Tahunan ")
}
func monthInterval(t time.Time) (firstDay, lastDay time.Time) {
	y, m, _ := t.Date()
	loc := t.Location()
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, loc)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, loc)
	return firstDay, lastDay
}
func GetDataDashboardReport() {
	functions.Logger().Info("Starting Scheduler Get Data Dashboard Report ")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file koneksi db")
	}
	db, err := sql.Open(os.Getenv("DB_CONNECTION"), ""+os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@tcp("+os.Getenv("DB_HOST")+":"+os.Getenv("DB_PORT")+")/"+os.Getenv("DB_DATABASE")+"")
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(0)
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(time.Nanosecond)

	//cTgl := jodaTime.Format("YYYYMMdd", time.Now())
	//cMonth := jodaTime.Format("MM", time.Now())
	cYear := jodaTime.Format("YYYY", time.Now())
	DBRe := os.Getenv("DB_DATABASE_RE")
	m, _ := strconv.Atoi(jodaTime.Format("MM", time.Now()))
	DD, _ := strconv.Atoi(jodaTime.Format("dd", time.Now()))
	z := m
	s := strconv.Itoa(m)
	d := ""
	if z < 10 {
		d = "0" + s
	} else {
		d = s
	}

	z = DD
	s = strconv.Itoa(DD)
	dX := ""
	if z < 10 {
		dX = "0" + s
	} else {
		dX = s
	}
	newdate := cYear + "-" + d + "-" + dX
	t, _ := time.Parse("2006-01-02", newdate)
	t = t.AddDate(0, 0, -1)
	cNewDate := jodaTime.Format("YYMMdd", t)
	cNewDate2 := jodaTime.Format("YYYYMMdd", t)
	t = t.AddDate(0, -1, 0)
	_, last := monthInterval(t)
	cNewDate3 := jodaTime.Format("YYYYMMdd", last)

	//REPORT NERACA
	sqlStatement := "SELECT kode_kantor from app_kode_kantor order by kode_kantor "
	rowsY, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}
	for rowsY.Next() {
		messageList := SelectDashboardReal{}
		err = rowsY.Scan(&messageList.KodeKantor)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		KodeKantor := messageList.KodeKantor

		///DELETE DULU DATA
		SendAPIDelete("delReportNeraca" + "/" + cNewDate + "/" + KodeKantor)
		_, err := db.Query("drop view IF EXISTS cdashboard_report_neraca ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err := db.Query(sqlStatementX)

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_neraca " +
				"as " +
				"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
				"" + cNewDate2 + " periode,a.type_perk " +
				"from " + DBRe + ".acc_" + cNewDate + " a,perkiraan b " +
				"where a.kode_perk=b.kode_perk " +
				"and (a.type_perk='HARTA') and a.kode_kantor='" + KodeKantor + "'" +
				"group by a.kode_perk having saldo<>0 order by a.kode_perk")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_neraca"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardReport{}
				err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
					&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				Amount := fmt.Sprintf("%.2f", messageList.Amount)
				Level := strconv.Itoa(messageList.Level)
				SendAPIPost("setReportNeraca" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
					messageList.KodeInduk + "/" + Level + "/" +
					messageList.GorD + "/" + messageList.NamaPerk +
					"/" + messageList.TypePerk + "/" + Amount + "/" + KodeKantor)

			}
		}

		///DELETE DULU DATA
		_, err = db.Query("drop view IF EXISTS cdashboard_report_neraca ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_neraca " +
				"as " +
				"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
				"" + cNewDate2 + " periode,a.type_perk " +
				"from " + DBRe + ".acc_" + cNewDate + "  a,perkiraan b " +
				"where a.kode_perk=b.kode_perk " +
				"and (a.type_perk='KEWAJIBAN') and a.kode_kantor='" + KodeKantor + "'" +
				"group by a.kode_perk having saldo<>0 order by a.kode_perk")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_neraca"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardReport{}
				err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
					&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				Amount := fmt.Sprintf("%.2f", messageList.Amount)
				Level := strconv.Itoa(messageList.Level)
				SendAPIPost("setReportNeraca" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
					messageList.KodeInduk + "/" + Level + "/" +
					messageList.GorD + "/" + messageList.NamaPerk +
					"/" + messageList.TypePerk + "/" + Amount + "/" + KodeKantor)

			}
		}

		///DELETE DULU DATA
		_, err = db.Query("drop view IF EXISTS cdashboard_report_neraca ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_neraca " +
				"as " +
				"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
				"" + cNewDate2 + " periode,a.type_perk " +
				"from " + DBRe + ".acc_" + cNewDate + "  a,perkiraan b " +
				"where a.kode_perk=b.kode_perk " +
				"and (a.type_perk='MODAL' or a.type_perk='LABA RUGI') and a.kode_kantor='" + KodeKantor + "'" +
				"group by a.kode_perk having saldo<>0 order by a.kode_perk")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_neraca"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardReport{}
				err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
					&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				Amount := fmt.Sprintf("%.2f", messageList.Amount)
				Level := strconv.Itoa(messageList.Level)
				SendAPIPost("setReportNeraca" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
					messageList.KodeInduk + "/" + Level + "/" +
					messageList.GorD + "/" + messageList.NamaPerk +
					"/" + messageList.TypePerk + "/" + Amount + "/" + KodeKantor)

			}
		}

		//REPORT LABARUGI
		///DELETE DULU DATA
		SendAPIDelete("delReportLabaRugi" + "/" + cNewDate + "/" + KodeKantor)

		_, err = db.Query("drop view IF EXISTS cdashboard_report_labarugi ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_labarugi " +
				"as " +
				"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
				"" + cNewDate2 + " periode,a.type_perk " +
				"from " + DBRe + ".acc_" + cNewDate + " a,perkiraan b " +
				"where a.kode_perk=b.kode_perk " +
				"and (a.type_perk='PENDAPATAN') and a.kode_kantor='" + KodeKantor + "' " +
				"group by a.kode_perk having saldo<>0 order by a.kode_perk")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_labarugi"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardReport{}
				err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
					&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				Amount := fmt.Sprintf("%.2f", messageList.Amount)
				Level := strconv.Itoa(messageList.Level)
				SendAPIPost("setReportLabaRugi" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
					messageList.KodeInduk + "/" + Level + "/" +
					messageList.GorD + "/" + messageList.NamaPerk +
					"/" + messageList.TypePerk + "/" + Amount + "/" + KodeKantor)

			}
		}
		//END OF REPORT LABARUGI

		//REPORT LABARUGI
		///DELETE DULU DATA

		_, err = db.Query("drop view IF EXISTS cdashboard_report_labarugi ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " where kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_labarugi " +
				"as " +
				"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
				"" + cNewDate2 + " periode,a.type_perk " +
				"from " + DBRe + ".acc_" + cNewDate + " a,perkiraan b " +
				"where a.kode_perk=b.kode_perk " +
				"and (a.type_perk='BIAYA' or a.type_perk='PAJAK') and a.kode_kantor='" + KodeKantor + "' " +
				"group by a.kode_perk having saldo<>0 order by a.kode_perk")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_labarugi"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rows.Next() {
				messageList := SelectDashboardReport{}
				err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
					&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}

				Amount := fmt.Sprintf("%.2f", messageList.Amount)
				Level := strconv.Itoa(messageList.Level)
				SendAPIPost("setReportLabaRugi" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
					messageList.KodeInduk + "/" + Level + "/" +
					messageList.GorD + "/" + messageList.NamaPerk +
					"/" + messageList.TypePerk + "/" + Amount + "/" + KodeKantor)

			}
		}
		//END OF REPORT LABARUGI

		//REPORT KREDIT
		///DELETE DULU DATA

		SendAPIDelete("delReportKredit" + "/" + cNewDate + "/" + KodeKantor)
		_, err = db.Query("drop view IF EXISTS cdashboard_report_rekap_kredit ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_rekap_kredit " +
				"as " +
				"Select a.kolektibilitas,a.my_kolek_number,sum(a.baki_debet) baki_debet,sum(if(a.baki_debet>0,1,0)) jml_rek," +
				"'" + cNewDate2 + "' periode " +
				"from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening and kredit.kode_kantor='" + KodeKantor + "' " +
				"group by a.my_kolek_number order by a.my_kolek_number")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select periode,my_kolek_number," +
				"if(kolektibilitas='L','Lancar'," +
				"if(kolektibilitas='DP','Dalam Perhatian Khusus'," +
				"if(kolektibilitas='DPK','Dalam Perhatian Khusus'," +
				"if(kolektibilitas='KL','Kurang Lancar'," +
				"if(kolektibilitas='D','Diragukan','Macet'))))) as kolektibilitas," +
				"baki_debet,jml_rek from cdashboard_report_rekap_kredit"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardReportKredit{}
				err = rows.Scan(&messageList.Period, &messageList.Kode, &messageList.MyKolekNumber, &messageList.BakiDebet, &messageList.JmlRek)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				//INSERT INTO DASHBOARD_CHART_TABUNGAN
				Amount := fmt.Sprintf("%.2f", messageList.BakiDebet)
				Jml := strconv.Itoa(messageList.JmlRek)
				SendAPIPost("setReportKredit" + "/" + messageList.Period + "/" + messageList.Kode + "/" +
					messageList.MyKolekNumber + "/" + Amount + "/" + Jml + "/" + KodeKantor)

			}
		}
		//END OF REPORT KREDIT

		//REPORT TABUNGAN
		///DELETE DULU DATA
		SendAPIDelete("delReportTabungan" + "/" + cNewDate + "/" + KodeKantor)

		_, err = db.Query("drop view IF EXISTS cdashboard_report_rekap_tabungan ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".tabung_" + cNewDate + " a,tabung where a.no_rekening=tabung.no_rekening and tabung.kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err = db.Query("create view cdashboard_report_rekap_tabungan " +
				"as " +
				"Select b.kode_produk,deskripsi_produk nama_produk,sum(a.saldo_akhir) saldo_akhir,count(a.no_rekening) jml_rek," +
				"'" + cNewDate2 + "' periode " +
				"from " + DBRe + ".tabung_" + cNewDate + " a,tabung b,tab_produk c " +
				"where a.no_rekening=b.no_rekening " +
				"AND b.kode_produk not in ('ABP','AKP','KPE') " +
				"and b.kode_produk=c.kode_produk and b.kode_kantor='" + KodeKantor + "'" +
				"group by b.kode_produk order by b.kode_produk")

			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select periode,kode_produk,nama_produk,saldo_akhir,jml_rek from cdashboard_report_rekap_tabungan"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			for rows.Next() {
				messageList := SelectDashboardReportKredit{}
				err = rows.Scan(&messageList.Period, &messageList.Kode, &messageList.MyKolekNumber, &messageList.BakiDebet, &messageList.JmlRek)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				Amount := fmt.Sprintf("%.2f", messageList.BakiDebet)
				Jml := strconv.Itoa(messageList.JmlRek)
				SendAPIPost("setReportTabungan" + "/" + messageList.Period + "/" + messageList.Kode + "/" +
					messageList.MyKolekNumber + "/" + Amount + "/" + Jml + "/" + KodeKantor)

			}
		}
		//END OF REPORT TABUNGAN

		//REPORT DEPOSITO
		///DELETE DULU DATA

		SendAPIDelete("delReportDeposito" + "/" + cNewDate + "/" + KodeKantor)

		_, err = db.Query("drop view IF EXISTS cdashboard_report_rekap_deposito ")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatementX = "select 1 from " + DBRe + ".deposito_" + cNewDate + " a,deposito where a.no_rekening=deposito.no_rekening and deposito.kode_kantor='" + KodeKantor + "' limit 1"
		rowsX, err = db.Query(sqlStatementX)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rowsX.Next() {
			_, err := db.Query("create view cdashboard_report_rekap_deposito " +
				"as " +
				"Select b.kode_produk,deskripsi_produk nama_produk,sum(a.saldo_akhir) saldo_akhir,count(a.no_rekening) jml_rek," +
				"'" + cNewDate2 + "' periode " +
				"from " + DBRe + ".deposito_" + cNewDate + " a,deposito b,dep_produk c " +
				"where a.no_rekening=b.no_rekening " +
				"and b.kode_produk=c.kode_produk and b.kode_kantor='" + KodeKantor + "' " +
				"group by b.kode_produk order by b.kode_produk")
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			sqlStatement := "select periode,kode_produk,nama_produk,saldo_akhir,jml_rek from cdashboard_report_rekap_deposito"
			rows, err := db.Query(sqlStatement)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			for rows.Next() {
				messageList := SelectDashboardReportKredit{}
				err = rows.Scan(&messageList.Period, &messageList.Kode, &messageList.MyKolekNumber, &messageList.BakiDebet, &messageList.JmlRek)
				if err != nil {
					functions.Logger().Error(err.Error())
					return
				}
				Amount := fmt.Sprintf("%.2f", messageList.BakiDebet)
				Jml := strconv.Itoa(messageList.JmlRek)
				SendAPIPost("setReportDeposito" + "/" + messageList.Period + "/" + messageList.Kode + "/" +
					messageList.MyKolekNumber + "/" + Amount + "/" + Jml + "/" + KodeKantor)

			}
		}
		//END OF REPORT DEPOSITO

	}
	///DELETE DULU DATA
	SendAPIDelete("delReportNeraca" + "/" + cNewDate + "/All")
	_, err = db.Query("drop view IF EXISTS cdashboard_report_neraca ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX := "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1"
	rowsX, err := db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}
	for rowsX.Next() {
		_, err = db.Query("create view cdashboard_report_neraca " +
			"as " +
			"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
			"" + cNewDate2 + " periode,a.type_perk " +
			"from " + DBRe + ".acc_" + cNewDate + "  a,perkiraan b " +
			"where a.kode_perk=b.kode_perk " +
			"and (a.type_perk='HARTA') " +
			"group by a.kode_perk having saldo<>0 order by a.kode_perk")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_neraca"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReport{}
			err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
				&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			Amount := fmt.Sprintf("%.2f", messageList.Amount)
			Level := strconv.Itoa(messageList.Level)
			SendAPIPost("setReportNeraca" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
				messageList.KodeInduk + "/" + Level + "/" +
				messageList.GorD + "/" + messageList.NamaPerk +
				"/" + messageList.TypePerk + "/" + Amount + "/All")

		}

	}

	///DELETE DULU DATA
	_, err = db.Query("drop view IF EXISTS cdashboard_report_neraca ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1 "
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}
	for rowsX.Next() {
		_, err = db.Query("create view cdashboard_report_neraca " +
			"as " +
			"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
			"" + cNewDate2 + " periode,a.type_perk " +
			"from " + DBRe + ".acc_" + cNewDate + "  a,perkiraan b " +
			"where a.kode_perk=b.kode_perk " +
			"and (a.type_perk='KEWAJIBAN') " +
			"group by a.kode_perk having saldo<>0 order by a.kode_perk")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_neraca"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReport{}
			err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
				&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			Amount := fmt.Sprintf("%.2f", messageList.Amount)
			Level := strconv.Itoa(messageList.Level)
			SendAPIPost("setReportNeraca" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
				messageList.KodeInduk + "/" + Level + "/" +
				messageList.GorD + "/" + messageList.NamaPerk +
				"/" + messageList.TypePerk + "/" + Amount + "/All")

		}

	}

	///DELETE DULU DATA
	_, err = db.Query("drop view IF EXISTS cdashboard_report_neraca ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1 "
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsX.Next() {
		_, err = db.Query("create view cdashboard_report_neraca " +
			"as " +
			"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
			"" + cNewDate2 + " periode,a.type_perk " +
			"from " + DBRe + ".acc_" + cNewDate + "  a,perkiraan b " +
			"where a.kode_perk=b.kode_perk " +
			"and (a.type_perk='MODAL' or a.type_perk='LABA RUGI') " +
			"group by a.kode_perk having saldo<>0 order by a.kode_perk")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_neraca"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReport{}
			err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
				&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			Amount := fmt.Sprintf("%.2f", messageList.Amount)
			Level := strconv.Itoa(messageList.Level)
			SendAPIPost("setReportNeraca" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
				messageList.KodeInduk + "/" + Level + "/" +
				messageList.GorD + "/" + messageList.NamaPerk +
				"/" + messageList.TypePerk + "/" + Amount + "/All")

		}

	}

	//REPORT LABARUGI
	///DELETE DULU DATA
	SendAPIDelete("delReportLabaRugi" + "/" + cNewDate + "/All")

	_, err = db.Query("drop view IF EXISTS cdashboard_report_labarugi ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1"
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsX.Next() {
		_, err := db.Query("create view cdashboard_report_labarugi " +
			"as " +
			"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
			"" + cNewDate2 + " periode,a.type_perk " +
			"from " + DBRe + ".acc_" + cNewDate + " a,perkiraan b " +
			"where a.kode_perk=b.kode_perk " +
			"and (a.type_perk='PENDAPATAN') " +
			"group by a.kode_perk having saldo<>0 order by a.kode_perk")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_labarugi"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReport{}
			err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
				&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			Amount := fmt.Sprintf("%.2f", messageList.Amount)
			Level := strconv.Itoa(messageList.Level)
			SendAPIPost("setReportLabaRugi" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
				messageList.KodeInduk + "/" + Level + "/" +
				messageList.GorD + "/" + messageList.NamaPerk +
				"/" + messageList.TypePerk + "/" + Amount + "/All")

		}

	}

	//END OF REPORT LABARUGI
	//REPORT LABARUGI
	///DELETE DULU DATA
	SendAPIDelete("delReportLabaRugi" + "/" + cNewDate + "/All")

	_, err = db.Query("drop view IF EXISTS cdashboard_report_labarugi ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".acc_" + cNewDate + " limit 1 "
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsX.Next() {
		_, err := db.Query("create view cdashboard_report_labarugi " +
			"as " +
			"Select a.kode_perk,a.kode_induk,a.level_perk,a.g_or_d,b.nama_perk,sum(a.saldo_akhir) saldo," +
			"" + cNewDate2 + " periode,a.type_perk " +
			"from " + DBRe + ".acc_" + cNewDate + " a,perkiraan b " +
			"where a.kode_perk=b.kode_perk " +
			"and (a.type_perk='BIAYA' or a.type_perk='PAJAK') " +
			"group by a.kode_perk having saldo<>0 order by a.kode_perk")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select kode_perk,kode_induk,level_perk,g_or_d,replace(replace(nama_perk,'%','Persen'),'/','\\\\') nama_perk,saldo,periode,type_perk from cdashboard_report_labarugi"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReport{}
			err = rows.Scan(&messageList.KodePerk, &messageList.KodeInduk, &messageList.Level, &messageList.GorD,
				&messageList.NamaPerk, &messageList.Amount, &messageList.Period, &messageList.TypePerk)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}

			Amount := fmt.Sprintf("%.2f", messageList.Amount)
			Level := strconv.Itoa(messageList.Level)
			SendAPIPost("setReportLabaRugi" + "/" + messageList.Period + "/" + messageList.KodePerk + "/" +
				messageList.KodeInduk + "/" + Level + "/" +
				messageList.GorD + "/" + messageList.NamaPerk +
				"/" + messageList.TypePerk + "/" + Amount + "/All")

		}

	}

	//END OF REPORT LABARUGI

	//REPORT KREDIT
	///DELETE DULU DATA

	SendAPIDelete("delReportKredit" + "/" + cNewDate + "/All")
	_, err = db.Query("drop view IF EXISTS cdashboard_report_rekap_kredit ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".kredit_" + cNewDate + " a limit 1"
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsX.Next() {
		_, err := db.Query("create view cdashboard_report_rekap_kredit " +
			"as " +
			"Select a.kolektibilitas,a.my_kolek_number,sum(a.baki_debet) baki_debet,count(a.no_rekening) jml_rek," +
			"'" + cNewDate2 + "' periode " +
			"from " + DBRe + ".kredit_" + cNewDate + " a,kredit where a.no_rekening=kredit.no_rekening " +
			"group by a.my_kolek_number order by a.my_kolek_number")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select periode,my_kolek_number," +
			"if(kolektibilitas='L','Lancar'," +
			"if(kolektibilitas='DP','Dalam Perhatian Khusus'," +
			"if(kolektibilitas='DPK','Dalam Perhatian Khusus'," +
			"if(kolektibilitas='KL','Kurang Lancar'," +
			"if(kolektibilitas='D','Diragukan','Macet'))))) as kolektibilitas," +
			"baki_debet,jml_rek from cdashboard_report_rekap_kredit"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReportKredit{}
			err = rows.Scan(&messageList.Period, &messageList.Kode, &messageList.MyKolekNumber, &messageList.BakiDebet, &messageList.JmlRek)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			//INSERT INTO DASHBOARD_CHART_TABUNGAN
			Amount := fmt.Sprintf("%.2f", messageList.BakiDebet)
			Jml := strconv.Itoa(messageList.JmlRek)
			SendAPIPost("setReportKredit" + "/" + messageList.Period + "/" + messageList.Kode + "/" +
				messageList.MyKolekNumber + "/" + Amount + "/" + Jml + "/All")

		}

	}

	//END OF REPORT KREDIT

	//REPORT TABUNGAN
	///DELETE DULU DATA
	SendAPIDelete("delReportTabungan" + "/" + cNewDate + "/All")

	_, err = db.Query("drop view IF EXISTS cdashboard_report_rekap_tabungan ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".tabung_" + cNewDate + " limit 1"
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsX.Next() {
		_, err = db.Query("create view cdashboard_report_rekap_tabungan " +
			"as " +
			"Select b.kode_produk,deskripsi_produk nama_produk,sum(a.saldo_akhir) saldo_akhir,count(a.no_rekening) jml_rek," +
			"'" + cNewDate2 + "' periode " +
			"from " + DBRe + ".tabung_" + cNewDate + " a,tabung b,tab_produk c " +
			"where a.no_rekening=b.no_rekening " +
			"AND b.kode_produk not in ('ABP','AKP','KPE') " +
			"and b.kode_produk=c.kode_produk " +
			"group by b.kode_produk order by b.kode_produk")
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select periode,kode_produk,nama_produk,saldo_akhir,jml_rek from cdashboard_report_rekap_tabungan"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReportKredit{}
			err = rows.Scan(&messageList.Period, &messageList.Kode, &messageList.MyKolekNumber, &messageList.BakiDebet, &messageList.JmlRek)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			Amount := fmt.Sprintf("%.2f", messageList.BakiDebet)
			Jml := strconv.Itoa(messageList.JmlRek)
			SendAPIPost("setReportTabungan" + "/" + messageList.Period + "/" + messageList.Kode + "/" +
				messageList.MyKolekNumber + "/" + Amount + "/" + Jml + "/All")

		}

	}

	//END OF REPORT TABUNGAN

	//REPORT DEPOSITO
	///DELETE DULU DATA

	SendAPIDelete("delReportDeposito" + "/" + cNewDate + "/All")

	_, err = db.Query("drop view IF EXISTS cdashboard_report_rekap_deposito ")
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatementX = "select 1 from " + DBRe + ".deposito_" + cNewDate + " limit 1 "
	rowsX, err = db.Query(sqlStatementX)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rowsX.Next() {
		_, err = db.Query("create view cdashboard_report_rekap_deposito " +
			"as " +
			"Select b.kode_produk,deskripsi_produk nama_produk,sum(a.saldo_akhir) saldo_akhir,count(a.no_rekening) jml_rek," +
			"'" + cNewDate2 + "' periode " +
			"from " + DBRe + ".deposito_" + cNewDate + " a,deposito b,dep_produk c " +
			"where a.no_rekening=b.no_rekening " +
			"and b.kode_produk=c.kode_produk " +
			"group by b.kode_produk order by b.kode_produk")

		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		sqlStatement := "select periode,kode_produk,nama_produk,saldo_akhir,jml_rek from cdashboard_report_rekap_deposito"
		rows, err := db.Query(sqlStatement)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}

		for rows.Next() {
			messageList := SelectDashboardReportKredit{}
			err = rows.Scan(&messageList.Period, &messageList.Kode, &messageList.MyKolekNumber, &messageList.BakiDebet, &messageList.JmlRek)
			if err != nil {
				functions.Logger().Error(err.Error())
				return
			}
			Amount := fmt.Sprintf("%.2f", messageList.BakiDebet)
			Jml := strconv.Itoa(messageList.JmlRek)
			SendAPIPost("setReportDeposito" + "/" + messageList.Period + "/" + messageList.Kode + "/" +
				messageList.MyKolekNumber + "/" + Amount + "/" + Jml + "/All")

		}

	}

	//END OF REPORT DEPOSITO

	//REPORT DEPOSITO
	///DELETE DULU DATA
	SendAPIDelete("delReportTks" + "/" + cNewDate3)
	SendAPIDelete("delRekapTks" + "/" + cNewDate3)

	_, err = db.Query("drop view IF EXISTS cdashboard_report_tks ")

	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	_, err = db.Query("create view cdashboard_report_tks " +
		"as " +
		"Select kode,huruf,keterangan,rasio,tgl_trans from tks_2016 " +
		"where tgl_trans='" + cNewDate3 + "'")

	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	sqlStatement = "select kode,huruf,keterangan,rasio,tgl_trans from cdashboard_report_tks"
	rows, err := db.Query(sqlStatement)
	if err != nil {
		functions.Logger().Error(err.Error())
		return
	}

	for rows.Next() {
		messageList := SelectDashboardReportTKS{}
		err = rows.Scan(&messageList.NoID, &messageList.SandiPos, &messageList.Deskripsi, &messageList.Jumlah, &messageList.Tanggal)
		if err != nil {
			functions.Logger().Error(err.Error())
			return
		}
		Amount := fmt.Sprintf("%.2f", messageList.Jumlah)
		SendAPIPost("setReportTKS" + "/" + messageList.NoID + "/" + messageList.NoID + "/" +
			"G/" + messageList.Deskripsi + "/" + Amount + "/" + messageList.Tanggal)

		SendAPIPost("setRekapTKS" + "/" + messageList.Deskripsi + "/" + Amount + "/" + messageList.Tanggal)
	}

	//END OF REPORT TKS

	defer db.Close()
	functions.Logger().Info("Successfully Get Data Dashboard Report ")
}

func SendAPIDelete(ApiCode string) int {

	functions.Logger().Info("Starting Send API Delete " + ApiCode)
	uri := os.Getenv("URL_API") + "/" + ApiCode
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, bytes.NewBufferString(""))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	functions.Logger().Info("Successfully Send API Delete " + ApiCode)
	return 1
}

func SendAPIPost(ApiCode string) int {

	functions.Logger().Info("Starting Send API Post " + ApiCode)
	uri := os.Getenv("URL_API") + "/" + ApiCode
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, uri, bytes.NewBufferString(""))

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	_, err = client.Do(req)
	functions.Logger().Info("Successfully Send API Post " + ApiCode)
	return 1
}
