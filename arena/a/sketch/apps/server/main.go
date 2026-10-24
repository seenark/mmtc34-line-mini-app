// Command server runs the app: `server` serves HTTP, `server seed` loads dev data.
package main

// main: slog JSON handler; cfg := server.ConfigFromEnv(); sqlDB := db.Open(ctx, cfg.DataDir);
// if os.Args[1] == "seed" { db.Seed; exit }  (seed branch added by ticket .3);
// h := server.New(server.Deps{cfg, sqlDB, line.Client{ChannelID: cfg.LineChannelID}, time.Now});
// http.Server{Addr: ":" + cfg.Port, Handler: h, ReadHeaderTimeout: 10s} with graceful
// shutdown on SIGTERM.
func main() {
	panic("not implemented")
}
