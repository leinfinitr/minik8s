import re
import sys
from collections import Counter

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            file_name = int(sys.argv[1])
            with open(file_name, 'r') as f:
                text = f.read()
            words = re.findall(r'\b\w+\b', text.lower())
            word_counts = Counter(words)
            dictionary = dict(word_counts)
            dictionary = dictionary.items()
            dictionary = sorted(dictionary, key=lambda x: x[1], reverse=True)
            print(dictionary[0][1] + " " + dictionary[0][0])
        except ValueError:
            print(0)