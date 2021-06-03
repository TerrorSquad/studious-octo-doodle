# bicsv

This is a CLI tool written in Go that will help
with [Bulk Image import for Magento 2](https://docs.magento.com/user-guide/system/data-import-product-images.html). The
tool will generate a valid CSV file that contains the names of images that should be imported.

### Input

path to the directory that contains the image files.

### Output

CSV file will be printed to STDOUT

### Supported image files are:

- jpg
- jpeg
- png

## Usage

```shell
bicsv ./product_images

bicsv ~/Images/Magento2/productImages > product_images_to_import.csv
```

## Installation

From binaries (Linux)
Download the archive from the Release page.

# License

This software is licensed under the Apache license.