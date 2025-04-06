## net basics

// Listen for TCP connections
net.Listen("tcp", ":8080")

// Accept a new connection
l.Accept()

// Read from a connection
conn.Read(buf)

// Write to a connection
conn.Write([]byte("..."))

// Echo data (copy from conn to conn)
io.Copy(dst, src)

// Dial out to a remote TCP host
net.Dial("tcp", "example.com:80")

## strings

// Split into lines
strings.Split(str, "\r\n")

// Split a line into parts
strings.SplitN(str, " ", 3)
strings.SplitN(str, ":", 2)

// Trim whitespace
strings.TrimSpace(str)

// Check prefix
strings.HasPrefix(str, "Host:")