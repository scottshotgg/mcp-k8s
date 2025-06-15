package kubernetes_tool

// TODO: make something like this but for kubectl logs in order to try out SSE

/*
	text/event-stream: Classic SSE format
	application/x-ndjson+stream: Line-delimited JSON
	text/plain+stream: Raw streaming text
*/

// server.RegisterResource(
//     "tools://build-log",
//     "Build Logs",
//     "Live build logs streamed via SSE",
//     "application/x-ndjson+stream",
//     func() (*mcp_golang.ResourceResponse, error) {
//         stream := mcp_golang.NewStreamedResource("tools://build-log", "application/x-ndjson+stream")

//         go func() {
//             defer stream.Close()
//             stream.WriteLine("Starting build...")
//             time.Sleep(1 * time.Second)
//             stream.WriteLine("Compiling...")
//             time.Sleep(1 * time.Second)
//             stream.WriteLine("Done.")
//         }()

//         return mcp_golang.NewResourceResponse(stream), nil
//     })
