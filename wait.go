package chromedp

import (
	"context"
	"sync"
)

// 这玩意有个弊端就是，如果检查的Action可能同时触发的话，返回的waitIdx在每次调用后不一定一样
// 所有请不要使用可能同时触发的Action
func WaitOneOf(waitIdx *int, actions ...Action) Action {
	if len(actions) == 0 {
		panic("actions cannot be empty")
	}
	if waitIdx == nil {
		panic("waitIdx cannot be nil")
	}

	return ActionFunc(func(ctx context.Context) error {
		wg := &sync.WaitGroup{}
		defer wg.Wait()

		ctx, cancel := context.WithCancel(ctx)
		// 因为用的一个ctx，所以也会cancel其他的
		// 这样的话，只要有一个返回了，就会忽略其他的了
		defer cancel()

		type ret struct {
			idx int
			err error
		}
		retC := make(chan ret)

		for idx := 0; idx < len(actions); idx++ {
			action := actions[idx]

			wg.Add(1)
			go func(idx int) {
				defer wg.Done()

				err := action.Do(ctx)
				select {
				case retC <- ret{
					idx: idx,
					err: err,
				}:
				case <-ctx.Done():
				}
			}(idx)
		}

		// 只要有一个返回了，就认作结束，然后上述的cancel会停掉所有的检查
		select {
		case r := <-retC:
			if r.err != nil {
				return r.err
			}

			*waitIdx = r.idx
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
}
