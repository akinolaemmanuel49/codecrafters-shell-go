git add .
echo "Enter commit message:"
read commit_message
git commit --allow-empty -m "$commit_message"
git push origin master
