package service

func (s Server) ApiRoutes() {
	s.Router.HandleFunc("/v0/control", s.NodeControlHandler)
	s.Router.HandleFunc("/v0/node", s.NodeStatsHandler)
	s.Router.HandleFunc("/v0/transactions", s.TransactionsHandler)
	s.Router.HandleFunc("/v0/model", s.ModelContextHandler)
	s.Router.HandleFunc("/v0/signal", s.SignalHandler)
}
