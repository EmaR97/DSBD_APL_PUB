import warnings

import dill
import pandas as pd
import seaborn as sns
from matplotlib import pyplot as plt

from UtilityPlot import plot_time_series, display_probability_surface
from time_series.gaussian_probability_estimation import calculate_mean_probability
from time_series.model_fitting import (apply_hp_filter_with_optimal_lambda, seasonal_decomposition,
                                       polynomial_fit_and_plot, polynomial_curve_function, seasonal_curve_fit_and_plot,
                                       seasonal_curve_function, complete_fit_and_plot,
                                       check_fitting_quality_and_print_metrics,
                                       print_error_distribution_and_return_stats, )

# Suppress all DeprecationWarnings
warnings.filterwarnings("ignore")


def main(data_csv, x_lower_limit, x_upper_limit, y_lower_bound, y_upper_bound):
    # Load the dataset

    df = pd.read_csv(data_csv)
    df['timestamp_sec'], df['value'] = df["X"], df["Y"]

    # Apply Hodrick-Prescott filter to decompose the time series into trend and cycle components
    cycle, trend, optimal_lambda = apply_hp_filter_with_optimal_lambda(time_series=df['value'],
                                                                       lambda_range=[1, 10, 50, 100, 200, 500, 0.5])

    # Plot the original time series, trend, and cycle components
    fig = plt.figure(figsize=(10, 12))
    plt.subplot(3, 1, 1)
    plot_time_series(df['timestamp_sec'], df['value'], label='Original Time Series', title='Original Time Series')
    plt.subplot(3, 1, 2)
    plot_time_series(df['timestamp_sec'], trend, label='Trend', title=f'Best Trend (lambda={optimal_lambda})',
                     linestyle='--', color='red')
    plt.subplot(3, 1, 3)
    plot_time_series(df['timestamp_sec'], cycle, label='Cycle', title=f'Best Cycle (lambda={optimal_lambda})',
                     linestyle='--', color='green')
    plt.tight_layout()
    plt.show()
    fig.savefig('images/original_time_series.png')  # Save the image with a name

    # Perform seasonal decomposition on the trend component
    result = seasonal_decomposition(time_series=trend, period_range=[10, 20, 50, 100, 200, 500, 1000, 2000])

    # Plot the decomposed components
    fig = plt.figure(figsize=(20, 20))
    plt.subplot(4, 1, 1)
    plot_time_series(df['timestamp_sec'], trend, label='Base Trend', title='Base Trend', linestyle='--', color='red')
    plt.subplot(4, 1, 2)
    plot_time_series(df['timestamp_sec'], result.trend, label='Trend', title='Trend', linestyle='--', color='blue')
    plt.subplot(4, 1, 3)
    plot_time_series(df['timestamp_sec'], result.seasonal, label='Seasonal', title='Seasonal', linestyle='--',
                     color='blue')
    plt.subplot(4, 1, 4)
    plot_time_series(df['timestamp_sec'], result.resid, label='Residual', title='Residual', linestyle='--',
                     color='purple')
    plt.tight_layout()
    plt.show()
    fig.savefig('images/decomposed_components.png')  # Save the image with a name

    # Fit a polynomial curve to the trend component
    df_trend = df.copy()
    df_trend["value"] = result.trend
    df_trend = df_trend.dropna()
    poly_coef_ = polynomial_fit_and_plot(f=polynomial_curve_function, x=df_trend["timestamp_sec"], y=df_trend["value"])

    # Plot the polynomial regression
    fig = plt.figure(figsize=(10, 10))
    plot_time_series(df_trend["timestamp_sec"], df_trend["value"], title='Base', label='Base', linestyle='-',
                     color='red')
    plot_time_series(df_trend["timestamp_sec"],
                     polynomial_curve_function(df_trend["timestamp_sec"], poly_coef_[0], poly_coef_[1],
                                               poly_coef_[2], ), title='Polynomial Regression',
                     label='Polynomial Regression', linestyle='--', color='blue')
    plt.show()
    fig.savefig('images/polynomial_regression.png')  # Save the image with a name

    # Fit a seasonal curve to the seasonal component
    period_coef_ = seasonal_curve_fit_and_plot(f=seasonal_curve_function, x=df['timestamp_sec'], y=result.seasonal)

    fig = plt.figure(figsize=(10, 10))
    plot_time_series(df['timestamp_sec'], result.seasonal, label='Base', title='Base', linestyle='-', color='red')
    plot_time_series(df['timestamp_sec'],
                     seasonal_curve_function(df['timestamp_sec'], period_coef_[0], period_coef_[1], period_coef_[2],
                                             period_coef_[3], ), label='Fitted', title='Curve Fit', linestyle='--',
                     color='blue')
    plt.show()
    fig.savefig('images/curve_fit.png')  # Save the image with a name

    # Complete the fitting process by combining polynomial and seasonal fits
    y_fitted_, fitting_error_, fitted_function_ = complete_fit_and_plot(df["timestamp_sec"], df["value"], poly_coef_,
                                                                        period_coef_)

    # Plot the base trend, fitted trend, and fitting error
    fig = plt.figure(figsize=(10, 10))
    fig.add_subplot(211)
    plot_time_series(df["timestamp_sec"], df["value"], title='Base', label='Base', linestyle='-', color='red')
    plot_time_series(df["timestamp_sec"], y_fitted_, title='Fitted', label='Fitted', linestyle='--', color='blue')
    fig.add_subplot(212)
    plot_time_series(df["timestamp_sec"], fitting_error_, title='Error', label='Error', linestyle='-', color='green')
    plt.show()
    fig.savefig('images/base_and_fitted_trend.png')  # Save the image with a name

    # Check the quality of fitting and print metrics
    check_fitting_quality_and_print_metrics(df["value"], y_fitted_)
    e_mean, e_std = print_error_distribution_and_return_stats(fitting_error_)

    # Visualize the distribution of errors
    fig = plt.figure(figsize=(20, 20))
    sns.histplot(fitting_error_, bins=30, kde=True, color='green')
    plt.title('Distribution of Errors')
    plt.xlabel('Error')
    plt.ylabel('Frequency')
    plt.show()
    fig.savefig('images/error_distribution.png')  # Save the image with a name

    serialized_function = dill.dumps(fitted_function_)
    deserialized_function = dill.loads(serialized_function)

    # Calculate the probability of y being out of bounds
    mean_probability, _, _, base_values = calculate_mean_probability(fitted_function_, e_std, x_lower_limit,
                                                                     x_upper_limit, threshold_u=y_upper_bound,
                                                                     threshold_l=y_lower_bound, num_points=1000)

    display_probability_surface(base_values, fitted_function_, e_std, threshold_u=y_upper_bound,
                                threshold_l=y_lower_bound)


if __name__ == "__main__":
    # main('data/data.csv', 1706547140, 1706549000, 2.06e+9, 2.08e+9)

    main('data/generated_data.csv', 0, 10, 0, 1)
