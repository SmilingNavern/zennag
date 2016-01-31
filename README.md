[![Build Status](https://travis-ci.org/SmilingNavern/zennag.svg?branch=master)](https://travis-ci.org/SmilingNavern/zennag)

# zennag

Simple HTTP-only monitoring written in golang with bolt database for store statistic.

## features
  * Simple monitoring almost without configuration
  * One yaml config
  * Bolt key/value database which requires no configuration at all
  * Static builded binary with no dependecies(thanks golang)
  * Alerts/Recovery with flap prevention(hard coded)

## todo
  * add alerter(telegram)
  * rewrite workers for WaitGroup
  * add agregation statistic
  * use cobra for cli(https://github.com/spf13/cobra)
  * add logging for errors
  * fix panic with running alerter on empty base
