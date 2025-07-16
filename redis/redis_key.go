package redis

// RedisCaptchaPrefix is the prefix of captcha key in redis
const RedisCaptchaPrefix = "CAPTCHA:"

// RedisTokenPrefix is the prefix of blacklist token key in redis
const RedisTokenPrefix = "BLACKLIST:TOKEN:"

// RedisTenantBlacklistPrefix is the prefix of tenant blacklist key in redis
const RedisTenantBlacklistPrefix = "BLACKLIST:TENANT:"

// RedisCasbinChannel is the channel of captcha key in redis
const RedisCasbinChannel = "/casbin"

// RedisApiPermissionCountPrefix is the prefix of api permission access times left in redis
const RedisApiPermissionCountPrefix = "API:PERMISSION:"

// RedisDataPermissionPrefix is the prefix of data permission in redis
const RedisDataPermissionPrefix = "DATAPERM:"

// RedisDynamicConfigurationPrefix is the prefix of dynamic configuration in redis
const RedisDynamicConfigurationPrefix = "CONFIGURATION:"
