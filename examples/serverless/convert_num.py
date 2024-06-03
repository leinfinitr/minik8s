# 将字符串中的每个字母转化为一个两位的数字
# a-z -> 01-26
# A-Z -> 27-52
import sys

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            string = sys.argv[1]
            result = ""
            for char in string:
                if char.islower():
                    result += str(ord(char) - ord('a') + 1).zfill(2)
                elif char.isupper():
                    result += str(ord(char) - ord('A') + 27).zfill(2)
                else:
                    result += char
        except ValueError:
            print(0)
