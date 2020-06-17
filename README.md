# go-application
Go application

[![Build Status](https://travis-ci.org/qumonintelligence/go-application.svg?branch=master)](https://travis-ci.org/qumonintelligence/go-application)

# start an application

```
func run1(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// done, exit loop now
			return
		default:
		}
	}
}

func run2(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// done, exit loop now
			return
		default:
		}
	}
}


func main() {
app := application.NewApplication(nil)
app.Start(run1)  // start a goroutine for run1
app.Start(run2)  // start a goroutine for run2
app.Background()  // wait until ctrl-c is pressed
}
```