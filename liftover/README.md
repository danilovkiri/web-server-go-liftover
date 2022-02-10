# 23andme Data Liftover hg38/hg19

## Description
The project can be used for positional liftover between the GRCh38 and GRCh37 builds. 

## Dependencies
1. [pyliftover](https://github.com/konstantint/pyliftover) (necessary contents are included)
2. UCSC liftover chain files (included), download is available for [hg19]([http://hgdownload.soe.ucsc.edu/goldenPath/hg19/liftOver/)
   and [hg38](http://hgdownload.soe.ucsc.edu/goldenPath/hg38/liftOver/)
6. Python 3+

## Usage
Run the `main.py` passing passing input, output text files and genome source build (either `hg38` or `hg19`).

```bash
python3 main.py <INPUT> <OUTPUT> <hg38|hg19>
```