import sys

if __name__ == "__main__":
    # 检查是否有命令行参数传入
    if len(sys.argv) != 2:
        print(0)
    else:
        try:
            dict = int(sys.argv[1])
            # 对输入的字典进行排序
            dict = dict.items()
            dict = sorted(dict, key=lambda x: x[1], reverse=True)
            # 输出最大的出现次数
            print(dict[0][1])
        except ValueError:
            print(0)