package main

func (t *TAdUgYybTypeInfoImpl) Select() (data []*TAdUgYybTypeInfoEntity, err error) {
	session := t.Connection.NewSession()
	_, err = session.Select("id,material,type").From("t_ad_ug_yyb_type_info").Load(&data)
	if err != nil {
		return
	}
	return
}
