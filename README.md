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

#### Server Options(use with `serve`

 - `-newsdir`: directory to serve newsfeed from
 - `-statsfile`: file to store the stats in, in json format
 - `-host`: host to serve news files on
 - `-port`: port to serve news files on
 - `-i2p`: serve news files directly to I2P using SAMv3

#### Builder Options(use with `build`

Not implemented yet

#### Signer Options(use with `sign`

Not implemented yet
