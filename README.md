# Meerkat
A CLI tool for fetching multiple public cameras composing a general global view.
Some web cameras only provide an static JPEG file that will refresh within a time frame.
This tool helps pull that JPEG's composing one MJPEG flow for each camera and
presenting them on a browser.


## How to install

You can go to the [releases page](https://github.com/eloylp/meerkat/releases) and download the lastest
binary for your architecture.

You can also try this one liner install:
```bash

```

## How to run 

If you want a description about the arguments, just invoke help:
```bash
meerkat -h
```

Here is an example of how to use this CLI:
```bash
meerkat -i 10 -l 0.0.0.0:3000 -u http://example.com/cam1.jpeg,http://example.com/cam2.jpeg
```

This will bring up a web server with the composed global camera view.