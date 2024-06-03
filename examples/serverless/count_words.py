import sys
from collections import Counter

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            words = int(sys.argv[1])
            word_counts = Counter(words)
            print(dict(word_counts))
        except ValueError:
            print(0)