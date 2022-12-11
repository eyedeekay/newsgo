newsgo
======

I2P News Server Tool/Library

Usage
-----

./newsgo -command $command -newsdir $news_directory -statsfile $news_stats_file

### Commands

 - serve: Serve newsfeeds from a directory
 - build: Build newsfeeds from XML(Not Implemented Yet)
 - sign: Sign newsfeeds with local keys(Not Implemented Yet)

### Options

Use these options to configure the software

#### Server Options(use with `serve`)

 - `-newsdir`: directory to serve newsfeed from
 - `-statsfile`: file to store the stats in, in json format
 - `-host`: host to serve news files on
 - `-port`: port to serve news files on
 - `-http`: serve news on host:port using HTTP
 - `-i2p`: serve news files directly to I2P using SAMv3

#### Builder Options(use with `build`)

 - `-newsfile`: entries to pass to news generator. If passed a directory, all `entries.html` files in the directory will be processed
 - `-blockfile`: block list file to pass to news generator
 - `-releasejson`: json file describing an update to pass to news generator
 - `-feedtitle`: title to use for the RSS feed to pass to news generator
 - `-feedsubtitle`: subtitle to use for the RSS feed to pass to news generator
 - `-feedsite`: site for the RSS feed to pass to news generator
 - `-feedmain`: Primary newsfeed for updates to pass to news generator
 - `-feedbackup`: Backup newsfeed for updates to pass to news generator
 - `-feeduri`: UUID to use for the RSS feed to pass to news generator

#### Signer Options(use with `sign`)

Not implemented yet
