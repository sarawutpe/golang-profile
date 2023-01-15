#!/bin/bash
echo "Setting up watches.";

# watch for any created, moved, or deleted image files
inotifywait -q -m -r --format '%e %w%f' -e close_write -e moved_from -e moved_to -e delete $1 \
| grep -i -E '\.(jpe?g|png)$' --line-buffered \
| while read operation path; do
  webp_path="$(sed 's/\.[^.]*$/.webp/' <<< "$path")";
  if [ $operation = "CLOSE_WRITE,CLOSE" ] || [ $operation = "MOVED_TO" ]; then  # if new file is created
     if [ $(grep -i '\.png$' <<< "$path") ]; then
       $(cwebp -quiet "$path" -o "$webp_path" && rm -f "$path");
     else
       $(cwebp -quiet "$path" -o "$webp_path" && rm -f "$path");
     fi;
       
  fi;
done;

