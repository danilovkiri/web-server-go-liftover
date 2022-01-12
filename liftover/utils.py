import sys
import subprocess
import platform
import warnings
warnings.filterwarnings('ignore')


def get_vartype(variation):
    if 'D' not in variation and 'I' not in variation:
        return 'SNV'
    else:
        return 'INDEL'


def catch_exception(pipe):
    if pipe.wait() != 1:
        pass
    else:
        print('### ERROR: Aborting')
        sys.exit(1)
    return None


def make_complement(string):
    complement = {'A':'T', 'T':'A', 'C':'G', 'G':'C', '/':'/', ':':':', 'N':'N', 'W':'W', 'S':'S', 'M':'M', 'K':'K',
                  'R':'R', 'Y':'Y', 'B':'B', 'D':'D', 'H':'H', 'V':'V', 'Z':'Z', '-':'-'}
    new_string = ''
    for i in string:
        new_string += complement[i.upper()]
    return new_string


def make_reverese_complement(string):
    complement = {'A':'T', 'T':'A', 'C':'G', 'G':'C', '/':'/', ':':':', 'N':'N', 'W':'W', 'S':'S', 'M':'M', 'K':'K',
                  'R':'R', 'Y':'Y', 'B':'B', 'D':'D', 'H':'H', 'V':'V', 'Z':'Z', '-':'-'}
    new_string = ''
    for i in string[::-1]:
        new_string += complement[i.upper()]
    return new_string


def execute_fasta_query(bashcommand):
    with subprocess.Popen(bashcommand, stdout=subprocess.PIPE, shell=True, executable='bash') as pipe:
        query = ''
        for line in pipe.stdout:
            query += line.decode(encoding='utf-8')
    catch_exception(pipe)
    return query


def query_dna(chr_id, start, end, project_dir, build, os_platform):
    if os_platform == 'Darwin':
        executable = 'twoBitToFa_macos'
    elif os_platform == 'Linux':
        executable = 'twoBitToFa_linux'
    else:
        print('### ERROR: Aborting')
        sys.exit(1)
    start = int(start) - 1
    if chr_id == 'MT':
        chr_id = 'M'
    if build == 'hg19':
        bashCommand = '{0}/{1} {0}/hg19.2bit -seq=chr{2} -start={3} -end={4} /dev/stdout'.format(project_dir, executable, chr_id, start, end)
    elif build == 'hg38':
        bashCommand = '{0}/{1} {0}/hg38.2bit -seq=chr{2} -start={3} -end={4} /dev/stdout'.format(project_dir, executable, chr_id, start, end)
    else:
        print('### ERROR: Aborting')
        sys.exit(1)
    query = execute_fasta_query(bashCommand)
    sequence = query.split('\n')[1].upper()
    return sequence


def check_complement(triplet1, triplet2):
    if make_complement(triplet1) == triplet2:
        return True
    elif make_reverese_complement(triplet1) == triplet2:
        return True
    else:
        return False


def check_system():
    os_platform = platform.system()
    if os_platform not in ['Darwin', 'Linux']:
        print('### ERROR: Unable to run all functions on non-Unix systems. Terminating.')
        sys.exit(1)
    return os_platform
