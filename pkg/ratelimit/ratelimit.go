package ratelimit

import (
	"strconv"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	limiter *redis_rate.Limiter
}

func New(redisClient *redis.Client) *Limiter {
	return &Limiter{
		limiter: redis_rate.NewLimiter(redisClient),
	}
}

func (rl *Limiter) Middleware(limit int, window time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.IP()
		res, err := rl.limiter.Allow(c.Context(), key, redis_rate.PerMinute(limit))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "rate limiter error",
			})
		}

		c.Set("X-Ratelimit-Limit", strconv.Itoa(limit))
		c.Set("X-Ratelimit-Remaining", strconv.Itoa(res.Remaining))
		c.Set("X-Ratelimit-Reset", strconv.FormatInt(time.Now().Add(res.ResetAfter).Unix(), 10))

		if res.Allowed == 0 {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests. Try again in " + res.RetryAfter.String(),
			})
		}

		return c.Next()
	}
}
