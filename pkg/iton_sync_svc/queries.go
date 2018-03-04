package iton_sync_svc

const (
	qSchema =
		`SELECT * from information_schema.tables where TABLE_TYPE='BASE TABLE'`

	qPersonList =
		`SELECT Идентификатор id, Имя firstName, Фамилия lastName, Отчество middleName, cardNumber, case when cast(ДатаУвольнения as date) <= cast(SYSDATETIME() as date) then 1 else 0 end
		 FROM СотрудникНабор
		`

	qsetPersonCard =
		`UPDATE СотрудникНабор SET cardNum=:1 WHERE Идентификатор=:2`

	qsetTimeSh =
		`MERGE INTO ТабельрабочегоВремени as Targ
		USING (SELECT CAST ($4 AS datetime) Дата, $5 Идентификатор) as Src
		ON (Src.Дата = Targ.Дата and Src.Идентификатор = Targ.Идентификатор)
		WHEN MATCHED THEN UPDATE
			set ВремяВхода = $1, ВремяВыхода= $2, СуммаВремя = $3
		WHEN NOT MATCHED THEN INSERT
				(Дата, Идентификатор, ВремяВхода, ВремяВыхода, СуммаВремя) values($4,$5, $1,$2,$3);
		`
)