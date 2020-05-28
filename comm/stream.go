package comm

import "io"

// ReadResponse holds the data just been read.
type ReadResponse struct {
	Buffer    []byte
	TotalRead int
}

// ReadStream reads data from io.Reader. Function onRead will be invoked when some data have been read into buffer.
func ReadStream(r io.Reader, bufferSize int, onRead func(resp *ReadResponse) error) error {
	buf := make([]byte, bufferSize)
	counter := &Counter{}
	resp := &ReadResponse{}
	for {
		err := counter.Add(r.Read(buf))
		if err != nil && err != io.EOF {
			return err
		}
		if counter.N() > 0 {
			resp.Buffer = buf[:counter.N()]
			resp.TotalRead = counter.Count()
			err := onRead(resp)
			if err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
	}
	return nil
}
