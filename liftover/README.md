# 23andme Data Liftover hg38/hg19

## Description
The project can be used for positional and optionally allelic liftover between the GRCh38 and GRCh37 builds. Positional
lifover may be accompanied by reverse-complement difference detection, which will trigger an allelic complementary change.
This option is time-consuming resulting in over x1000 speed reduction.

## Dependencies
1. [pyliftover](https://github.com/konstantint/pyliftover) (necessary contents are included)
2. UCSC liftover chain files (included), download is available for [hg19]([http://hgdownload.soe.ucsc.edu/goldenPath/hg19/liftOver/)
   and [hg38](http://hgdownload.soe.ucsc.edu/goldenPath/hg38/liftOver/)
3. `twoBitTwoFa_linux` and `twoBitTwoFa_macos` shell scripts, available from [here](http://hgdownload.soe.ucsc.edu/admin/exe/)),
   should be accessible from main project directory and have the above-specified name (might be necessary to change file
   permissions to run by `sudo chmod -R 777 <file>`); included in the project
4. `hg38.2bit` genome sequence file for fast access (can be downloaded from UCSC repository
   [here](http://hgdownload.soe.ucsc.edu/goldenPath/hg38/bigZips/)), should be accessible from main project directory and
   have the above-specified name
5. `hg19.2bit` genome sequence file for fast access (can be downloaded from UCSC repository
   [here](http://hgdownload.soe.ucsc.edu/goldenPath/hg19/bigZips/)), should be accessible from main project directory and
   have the above-specified name
6. Python 3+

## Usage
Run the `main.py` passing `--input`, `--output` and `--source` arguments, where the `--source` denotes the genome build
of the input file. You can optionally use `--fast` argument with the value `True` for fast conversion (disables allelic
check).

```bash
python3 main.py --input input.txt --output output.txt --source hg38 --fast True
```