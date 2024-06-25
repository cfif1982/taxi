package application

// QUESTION: у меня ctxKey и const KeyDriverID встречается в двух хэндлерах: drivers/handler.go и здесь.
// Как их лучше описать в одном месте? или оставить так?
// создаю свой тип для ключей контекста. Нужно хранить id авторизованного водителя
type ctxKey string

const KeyDriverID ctxKey = "driver_id" //  ключ в контексте для поля driver_id
