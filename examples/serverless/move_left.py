# 将一个字符串中的每个字母向左移动7位，其他字符不变
import sys

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            string = sys.argv[1]
            result = ""
            for c in string:
                if c.isalpha():
                    if c.islower():
                        result += chr((ord(c) - ord('a') - 7) % 26 + ord('a'))
                    else:
                        result += chr((ord(c) - ord('A') - 7) % 26 + ord('A'))
                else:
                    result += c
        except ValueError:
            print(0)
