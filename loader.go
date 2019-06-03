package randgenloader

type RandgenLoader struct {

}

func (*RandgenLoader) Init()  {

}

func (*RandgenLoader) LoadData(zzFile string) (sqls []string)  {
	return nil
}

func (*RandgenLoader) Query(yyFile string) (sqls []string)  {
	return nil
}

func (*RandgenLoader) Compare() (comment string, consistent bool)  {
	return "", true
}
