package middleware

import (
	"fmt"
	"github.com/plally/dgcommand"
	"math"
	"time"
)

func Coooldown(decayRate time.Duration, amount uint) dgcommand.MiddlewareFunc {
	return func(h dgcommand.HandlerFunc) dgcommand.HandlerFunc {
		cool := cooldownContainer{decayRate: decayRate, amount: amount, cooldowns: make(map[string]*cooldown)}

		return func(ctx dgcommand.CommandContext) {

			cooldownInstance, ok := cool.cooldowns[ctx.Message.Author.ID]
			if !ok {
				cooldownInstance = &cooldown{}
				cool.cooldowns[ctx.Message.Author.ID] = cooldownInstance
			}
			since := time.Since(cooldownInstance.lastRun)
			subAmount := uint(since / decayRate)

			if subAmount > cooldownInstance.amount {
				subAmount = cooldownInstance.amount
			}
			cooldownInstance.amount = cooldownInstance.amount - subAmount

			fmt.Println(cooldownInstance.amount, uint(math.Max(0, float64(since/decayRate))), since)
			if cooldownInstance.amount > cool.amount {
				ctx.Reply("You're doing that too much")
				return
			}

			cooldownInstance.lastRun = time.Now()
			cooldownInstance.amount++
			h(ctx)

		}
	}
}

type cooldown struct {
	lastRun time.Time
	amount  uint
}

type cooldownContainer struct {
	decayRate time.Duration
	amount    uint
	cooldowns map[string]*cooldown
}
