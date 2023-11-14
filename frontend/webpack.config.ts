import HTMLPlugin from 'html-webpack-plugin'
import ExtractPlugin from 'mini-css-extract-plugin'
import type { Configuration } from 'webpack'

const configuration: Configuration = {
  context: __dirname,
  resolve: {
    extensions: ['.ts', '.tsx', '.json', '.js'],
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        loader: require.resolve('ts-loader'),
      },
      {
        test: /\.css$/,
        use: [ExtractPlugin.loader, require.resolve('css-loader')],
      },
    ],
  },
  plugins: [
    new HTMLPlugin({
      title: 'eSIM LPA Agent',
      favicon: require.resolve('./e-sim.png'),
    }),
    new ExtractPlugin(),
  ],
}

export default configuration
