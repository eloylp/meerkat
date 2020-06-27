# Meerkat
A CLI tool for fetching multiple public cameras composing a general global view.
Some web cameras only provide a static JPEG file that will refresh periodically.
This tool helps to pull that JPEG's and composing one MJPEG flow for each camera
presenting them on a browser.

## How to install

You can go to the [releases page](https://github.com/eloylp/meerkat/releases) and download the latest
binary for your architecture.

You can also try this one liner install:
```bash
sudo curl -L "https://github.com/eloylp/meerkat/releases/download/v1.0.1/meerkat_1.0.1_Linux_x86_64" \
-o /usr/local/bin/meerkat \
&& sudo chmod +x /usr/local/bin/meerkat
```

## How to run 

If you want a description about the arguments, just invoke help:
```bash
meerkat -h
```

Here is a full example of how to use this CLI:
```bash
meerkat -i 10 -l 0.0.0.0:3000 -u http://example.com/cam1.jpeg,http://example.com/cam2.jpeg
```

This will bring up a web server with the composed global camera view.
Visit localhost:3000 in your browser to see results.