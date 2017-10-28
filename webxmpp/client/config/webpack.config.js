const packageJSON = require('../package.json');
const path = require('path');
const webpack = require('webpack');
const WebpackStrip = require('strip-loader');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const ExtractTextPlugin = require("extract-text-webpack-plugin");

const NODE_ENV = process.env.NODE_ENV || 'development';
const SERVER_RENDER = process.env.SERVER_RENDER === 'true'; // TODO
const PRODUCTION_BUILD = NODE_ENV === 'production';

const PATHS = {
  build: path.join(__dirname, '..', 'dist')
};

const extractSass = new ExtractTextPlugin({
    filename: "[name]-[contenthash].css",
    allChunks: true,
    disable: !PRODUCTION_BUILD
});

const plugins = [
  new HtmlWebpackPlugin({
    title: 'GoXMPP server',
    template: 'app/index.html',
    chunks: ['home-index'],
    filename: 'index.html'
  }),
  new webpack.DefinePlugin({
    'process.env': {
      'NODE_ENV': JSON.stringify(process.env.NODE_ENV)
    }
  }),
  extractSass
];

const standardModuleLoaders = {
    loaders: [
        {
            enforce: 'pre',
            test: /\.js|.jsx$/,
            loader: 'eslint-loader',
            exclude: /node_modules/,
            options: {
                emitWarning: true
            }
        },
        {
            test: /\.js|.jsx?$/,
            loader: 'babel-loader',
            exclude: !PRODUCTION_BUILD ? /node_modules/ : [
              /node_modules\/babel-/m,
              /node_modules\/core-js\//m,
              /node_modules\/regenerator-runtime\//m
            ],
            query: {
                presets: ['es2015', 'react', "stage-0"],
                plugins: ['transform-runtime', 'transform-decorators-legacy', 'transform-object-assign']
            }
        },
        {
            test: /\.scss|.css$/,
            use: extractSass.extract({
              fallback: 'style-loader',
              use: [
                { loader: 'css-loader', options: { minimize: PRODUCTION_BUILD, importLoaders: 1 } },
                'resolve-url-loader',
                'sass-loader?sourceMap'
              ]
            })
        },
        {
            test: /\.(woff|woff2|eot|ttf)$/,
            loader: 'file-loader',
            options: {
              outputPath: 'static/fonts/',
              name: '[name]-[hash].[ext]',
              emitFile: true
            }
        },
        {
            test: /\.(png|jpg|svg|jpeg|gif)$/,
            loader: 'file-loader',
            options: {
              outputPath: 'static/images/',
              name: '[name]-[hash].[ext]',
              emitFile: true
            }
        }
    ]
};

if (PRODUCTION_BUILD) {
  console.log("this is PRODUCTION_BUILD");
  standardModuleLoaders.loaders.push({
      test: /\.js|.jsx$/,
      // list of functions that should be removed
      loader: WebpackStrip.loader('console.log')
  });
  plugins.push(
    new webpack.optimize.ModuleConcatenationPlugin(),
    new webpack.optimize.UglifyJsPlugin({
      compress: {
        comparisons: true,
        conditionals: true,
        dead_code: true,
        drop_console: !SERVER_RENDER, // Keep server logs
        drop_debugger: true,
        evaluate: true,
        if_return: true,
        join_vars: true,
        screw_ie8: true,
        sequences: true,
        unused: true,
        warnings: false,
      },
      output: {
        comments: false,
      },
    })
  );
} else {
  plugins.push(
    new webpack.NamedModulesPlugin(),
    new webpack.optimize.OccurrenceOrderPlugin(),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.NoEmitOnErrorsPlugin()
  );
}

module.exports = [
    {
        entry: {
          'home-index': './app/bundles/home-index.js',
        },
        resolve: {
            extensions: ['.js', '.jsx']
        },
        module: standardModuleLoaders,
        output: {
            path: PATHS.build,
            publicPath: '/',
            filename: '[name]-[hash].js'
        },
        plugins: plugins,
        devServer: {
          historyApiFallback: {
            rewrites: [
              { from: /^\/$/, to: '/index.html' },
              { from: /./, to: '/index.html' },
            ],
          },
          port: 3000,
          publicPath: `http://localhost:3000/`,
          proxy: {
              '/rest': {
                  target: '/api',
                  // target: 'http://localhost:8080',
                  secure: false,
                  changeOrigin: true,
                  headers: {
                      "User-Agent": "Webpack-dev-server"
                  }
              }
          }
        }
    }
];
