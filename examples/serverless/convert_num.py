# 将字符串根据 convert_char 的逆规则转换为数字

# convert_char 规则：
# 若数字为 00，则转换为_
# 若数字介于 01-26，则转换为 a-z
# 若数字介于 27-52，则转换为 A-Z
# 若数字大于 53，则不做处理
import sys

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            string = sys.argv[1]
            result = ""
            for c in string:
                if c == "_":
                    result += "00"
                elif c.islower():
                    result += str(ord(c) - ord('a') + 1).zfill(2)
                elif c.isupper():
                    result += str(ord(c) - ord('A') + 27).zfill(2)
                else:
                    result += c
            print(result)
        except ValueError:
            print(0)
