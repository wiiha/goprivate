cd webfront
npm run build
echo "success: compiling webfront"
cd ..
echo "copying new webfront"
cp -r webfront/dist/ server/webfront/
echo "success: updating webfront"
echo "You can now compile goprivate"
