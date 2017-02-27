var ExtractTextPlugin = require("extract-text-webpack-plugin");
var path = require("path");

module.exports = {
    entry: [
        "webpack-dev-server/client?http://localhost:8080",
        "webpack/hot/only-dev-server",
        "./index.js"
    ],
    module: {
        loaders: [
            {
                test: /\.scss$/,
                loader: ExtractTextPlugin.extract(
                    "style",
                    "css!sass"
                )
            },
            {
                test: /\.(eot|svg|ttf|woff|woff2)$/,
                loader: 'file?name=public/fonts/[name].[ext]'
            }
        ]
    },
    resolve: {
        extensions: ["", ".js"],
        alias: {
            fonts: path.join(__dirname + "./../third_party/node/fonts"),
        },
    },
    resolveUrlLoader: {
        relative: '../third_party/node/fonts',
    },
    output: {
        path: __dirname + "/dist",
        publicPath: "./",
        filename: "bundle.js"
    },
    devServer: {
        contentBase: "./dist"
    },
    plugins: [
        new ExtractTextPlugin("styles.css")
    ]
};
