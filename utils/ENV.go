package utils

import "time"

const SECRET = "hmacSampleSecret"
const ACCESS_TOKEN_DURATION = time.Hour * 24
const REFRESH_TOKEN_DURATION = time.Hour * 24 * 90