package repoerrors

import (
	"fmt"
)

var (
	ErrorInsertShortLink        = fmt.Errorf("ошибка вставки короткой ссылки: ")
	ErrorSelectExistedShortLink = fmt.Errorf("ошибка при выборке существующей короткой ссылки: ")
	ErrorShortLinkAlreadyInDB   = fmt.Errorf("короткая ссылка уже в БД: ")
	ErrorSQLInternal            = fmt.Errorf("внутренняя ошибка базы данных: ")
	ErrorShortLinkNotFound      = fmt.Errorf("короткая ссылка не найдена: ")
	ErrorPingDB                 = fmt.Errorf("ошибка пинга базы данных: ")
	ErrorTableCreate            = fmt.Errorf("ошибка создания таблицы: ")
	ErrorConnectingDB           = fmt.Errorf("ошибка подключения к БД: ")
	ErrorSelectShortLinks       = fmt.Errorf("ошибка выборки коротких ссылок для пользователя ")
)
