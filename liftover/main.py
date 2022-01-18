#!/usr/bin/python3

import csv
import os
from pyliftover.liftover import LiftOver
import sys
import datetime
import numpy as np
import utils as utils
import warnings
warnings.filterwarnings('ignore')


def main(file_in, file_out, source_build):
    os_platform = utils.check_system()
    project_dir = os.path.abspath(os.path.join(os.path.abspath(__file__), '../'))
    if source_build == 'hg19':
        lo = LiftOver('{0}/hg19ToHg38.over.chain.gz'.format(project_dir))
    elif source_build == 'hg38':
        lo = LiftOver('{0}/hg38ToHg19.over.chain.gz'.format(project_dir))
    else:
        print('### ERROR: Incorrect source genome build')
        sys.exit(1)
    with open(file_in, 'r') as f_in, open(file_out, 'w+') as f_out:
        f_out.write('\t'.join(map(str, ['# rsid', 'chromosome', 'position', 'genotype'])) + '\n')
        reader = csv.reader(f_in, delimiter='\t')
        count = 1
        subtime0 = datetime.datetime.now()
        subtime1 = datetime.datetime.now()
        for i in reader:
            if count % 100000 == 0:
                subtime2 = datetime.datetime.now()
                diftime = subtime2 - subtime1
                difsec = diftime.seconds
                diftime0 = subtime2 - subtime0
                difmin0 = np.round(divmod(diftime0.seconds, 60)[0] + (divmod(diftime0.seconds, 60)[1] / 60), 2)
                print('### INFO: Progress: {0}\tLast 100000: {1} sec\tElapsed: {2} min'.format(count, difsec, difmin0))
                subtime1 = datetime.datetime.now()
            count += 1
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
                    vartype = utils.get_vartype(gt)
                    if vartype == 'INDEL':
                        f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, gt])) + '\n')
                    elif vartype == 'SNV':
                        if fast_mode is True:
                            f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, gt])) + '\n')
                        else:
                            if source_build == 'hg19':
                                seq_hg19 = utils.query_dna(chr_id, pos-4, pos+4, project_dir, 'hg19', os_platform)
                                seq_hg38 = utils.query_dna(chr_id, pos_lo-4, pos_lo+4, project_dir, 'hg38', os_platform)
                                allele_hg19 = seq_hg19[4]
                                allele_hg38 = seq_hg38[4]
                                trigger = utils.check_complement(seq_hg19, seq_hg38)
                                if trigger is False:
                                    f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, gt])) + '\n')
                                elif trigger is True:
                                    f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, utils.make_complement(gt)])) + '\n')
                            elif source_build == 'hg38':
                                seq_hg19 = utils.query_dna(chr_id, pos_lo-4, pos_lo+4, project_dir, 'hg19', os_platform)
                                seq_hg38 = utils.query_dna(chr_id, pos-4, pos+4, project_dir, 'hg38', os_platform)
                                allele_hg19 = seq_hg19[4]
                                allele_hg38 = seq_hg38[4]
                                trigger = utils.check_complement(seq_hg19, seq_hg38)
                                if trigger is False:
                                    f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, gt])) + '\n')
                                elif trigger is True:
                                    f_out.write('\t'.join(map(str, [snp_id, chr_id, pos_lo, utils.make_complement(gt)])) + '\n')
    f_in.close()
    f_out.close()
    return None


if __name__ == '__main__':
    file_in = str(sys.argv[1])
    file_out = str(sys.argv[2])
    source_build = str(sys.argv[3])
    fast_mode = True
    main(file_in, file_out, source_build)
