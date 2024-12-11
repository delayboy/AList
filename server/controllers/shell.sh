if echo "准备下载脚本"; then

%s

echo "一共${#my_filenames[@]}个文件，开始下载..."
# 使用for循环按索引遍历数组并打印每个元素
for ((i = 0; i < ${#my_dirs[@]}; i++))
do
  if ! test -d "${my_dirs[i]}"; then
    echo "$i 创建文件夹: ${my_dirs[i]}"
    mkdir -p "${my_dirs[i]}"
  else
    echo "$i 文件夹已经存在: ${my_dirs[i]}"
  fi
done

for ((i = 0; i < ${#my_filenames[@]}; i++))
do
  if ! test -f "${my_filenames[i]}"; then
    echo "($((1+i))/${#my_filenames[@]}) 文件不存在,下载文件: ${my_filenames[i]}"
    wget  --user-agent="pan.baidu.com" "${my_links[i]}" -O "${my_filenames[i]}"
  else
    echo "($((1+i))/${#my_filenames[@]}) 文件已存在跳过下载: ${my_filenames[i]}"
  fi
done
fi
my_dirs=( "./笔记工具" )
my_links=( "https://d.pcs.baidu.com/file/fde48abedt1cc548ac1c79863cd6bc7d?fid=3496751938-250528-204127907501754&rt=pr&sign=FDtAERV-DCb740ccc5511e5e8fedcff06b081203-1r3uDOFWYWS9%2Fu65jHU6IgFRm%2Bs%3D&expires=8h&chkbd=0&chkv=3&dp-logid=3252328059702956379&dp-callid=0&dstime=1691304737&r=656644399&origin_appid=28156763&file_type=0&access_token=121.7f7501691eadcdf3fcf8b2a33f1342b1.YBvT-C1PqqVsQOvWt2txPE4SulMK73rVOGA_tqw.5Efweg" \
 "https://d.pcs.baidu.com/file/7b445b4a0l9a984d08de1b0bb81c5725?fid=3496751938-250528-378366274910675&rt=pr&sign=FDtAERV-DCb740ccc5511e5e8fedcff06b081203-RXRMrCSjyr14XD2PjDU8Mv2xeNI%3D&expires=8h&chkbd=0&chkv=3&dp-logid=3251645150556895108&dp-callid=0&dstime=1691304738&r=617024059&origin_appid=28156763&file_type=0&access_token=121.7f7501691eadcdf3fcf8b2a33f1342b1.YBvT-C1PqqVsQOvWt2txPE4SulMK73rVOGA_tqw.5Efweg" \
 "https://d.pcs.baidu.com/file/805d889a2t362e8b9cce91c1ffd0e2cc?fid=3496751938-250528-965705629871716&rt=pr&sign=FDtAERV-DCb740ccc5511e5e8fedcff06b081203-dtBv7sEp1YGyUVGPl7hFgZqnJkA%3D&expires=8h&chkbd=0&chkv=3&dp-logid=3251863771368379431&dp-callid=0&dstime=1691304738&r=617024059&origin_appid=28156763&file_type=0&access_token=121.7f7501691eadcdf3fcf8b2a33f1342b1.YBvT-C1PqqVsQOvWt2txPE4SulMK73rVOGA_tqw.5Efweg" \
 "https://d.pcs.baidu.com/file/06802589fhdf8bd1e3ce4b67ef3a40b8?fid=3496751938-250528-743033928872410&rt=pr&sign=FDtAERV-DCb740ccc5511e5e8fedcff06b081203-bMiRtN1euif%2FAGpkOQ%2F3%2BIgtXks%3D&expires=8h&chkbd=0&chkv=3&dp-logid=3252085189756250531&dp-callid=0&dstime=1691304738&r=617024059&origin_appid=28156763&file_type=0&access_token=121.7f7501691eadcdf3fcf8b2a33f1342b1.YBvT-C1PqqVsQOvWt2txPE4SulMK73rVOGA_tqw.5Efweg" \
)
my_filenames=( "./笔记工具/frame.js" \
 "./笔记工具/MyFileEditor.jar" \
 "./笔记工具/ReadMe.md" \
 "./笔记工具/TyporaUpload.py" \
)
