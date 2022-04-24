package ports

type HttpHandler interface {
	SetupRoutes()
	Run(addres string)
}
