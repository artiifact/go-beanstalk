package mock

import "github.com/IvanLutokhin/go-beanstalk"

func NewClient(in, out []string) *beanstalk.Client {
	return beanstalk.NewClient(NewConn(in, out))
}
