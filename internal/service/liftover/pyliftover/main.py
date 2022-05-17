#!/usr/bin/python3

import csv
import os
from pyliftover.liftover import LiftOver
import sys
import warnings
warnings.filterwarnings('ignore')


def main(file_in, file_out, source_build):
    project_dir = os.path.abspath(os.path.join(os.path.abspath(__file__), '../'))
    if source_build == 'hg19':
        lo = LiftOver('{0}/hg19ToHg38.over.chain.gz'.format(project_dir))
    elif source_build == 'hg38':
        lo = LiftOver('{0}/hg38ToHg19.over.chain.gz'.format(project_dir))
    else:
        raise ValueError("Incorrect source genome build")
    with open(file_in, 'r') as f_in, open(file_out, 'w+') as f_out:
        f_out.write('\t'.join(map(str, ['# rsid', 'chromosome', 'position', 'genotype'])) + '\n')
        reader = csv.reader(f_in, delimiter='\t')
        for i in reader:
            if i[0][0] == '#':
                continue
            else:
                snp_id = str(i[0])
                chr_id = str(i[1])
                pos = int(i[2])
                gt = str(i[3])
                if chr_id == 'MT':
                    f_out.write('\t'.join(map(str, [snp_id, chr_id, pos, gt])) + '\n')
                else:
                    if source_build == 'hg19':
                        try:
                            pos_lo = int(list(lo.convert_coordinate('chr{0}'.format(chr_id), pos - 1, '+')[0])[1]) + 1
                        except:
                            continue
                    elif source_build == 'hg38':
                        try:
                            pos_lo = int(list(lo.convert_coordinate('chr{0}'.format(chr_id), pos - 1, '+')[0])[1]) + 1
                        except:
                            continue
                    f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, gt])) + '\n')
    f_in.close()
    f_out.close()
    return None


if __name__ == '__main__':
    file_in = str(sys.argv[1])
    file_out = str(sys.argv[2])
    source_build = str(sys.argv[3])
    main(file_in, file_out, source_build)
