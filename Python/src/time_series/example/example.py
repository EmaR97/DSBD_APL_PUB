import dill
import numpy as np
import pandas as pd
import seaborn as sns
from matplotlib import pyplot as plt

# Import utility functions for plotting
from UtilityPlot import plot_time_series, plot_surface
# Import functions for time series analysis
from time_series.gaussian_probability_estimation import generate_surface_function, calculate_probability
from time_series.model_fitting import (apply_hp_filter_with_optimal_lambda, seasonal_decomposition,
                                       polynomial_fit_and_plot, polynomial_curve_function, seasonal_curve_fit_and_plot,
                                       seasonal_curve_function, complete_fit_and_plot,
                                       check_fitting_quality_and_print_metrics,
                                       print_error_distribution_and_return_stats, )

# Load the dataset
df = pd.read_csv('data.csv')
x_lower_limit, x_upper_limit = 547140, 549000
y_lower_bound = 2.06e+9
# df = pd.read_csv('generated_data.csv')
# # Define the range for probability calculation
# x_lower_limit, x_upper_limit = 0, 10
# y_lower_bound = 1

df['timestamp_sec'], df['value'] = df["X"], df["Y"]

# Apply Hodrick-Prescott filter to decompose the time series into trend and cycle components
cycle, trend, optimal_lambda = apply_hp_filter_with_optimal_lambda(time_series=df['value'],
                                                                   lambda_range=[1, 10, 50, 100, 200, 500, 0.5])

# Plot the original time series, trend, and cycle components
plt.figure(figsize=(10, 12))
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

# Perform seasonal decomposition on the trend component
result = seasonal_decomposition(time_series=trend, period_range=[10, 20, 50, 100, 200, 500, 1000, 2000])

# Plot the decomposed components
plt.figure(figsize=(20, 20))
plt.subplot(4, 1, 1)
plot_time_series(df['timestamp_sec'], trend, label='Base Trend', title='Base Trend', linestyle='--', color='red')
plt.subplot(4, 1, 2)
plot_time_series(df['timestamp_sec'], result.trend, label='Trend', title='Trend', linestyle='--', color='blue')
plt.subplot(4, 1, 3)
plot_time_series(df['timestamp_sec'], result.seasonal, label='Seasonal', title='Seasonal', linestyle='--', color='blue')
plt.subplot(4, 1, 4)
plot_time_series(df['timestamp_sec'], result.resid, label='Residual', title='Residual', linestyle='--', color='purple')
plt.tight_layout()
plt.show()

# Fit a polynomial curve to the trend component
df_trend = df.copy()
df_trend["value"] = result.trend
df_trend = df_trend.dropna()
poly_coef_ = polynomial_fit_and_plot(f=polynomial_curve_function, x=df_trend["timestamp_sec"], y=df_trend["value"])

# Plot the polynomial regression
plt.figure(figsize=(20, 20))
plot_time_series(df_trend["timestamp_sec"], df_trend["value"], title='Base', label='Base', linestyle='-', color='red')
plot_time_series(df_trend["timestamp_sec"],
                 polynomial_curve_function(df_trend["timestamp_sec"], poly_coef_[0], poly_coef_[1], poly_coef_[2], ),
                 title='Polynomial Regression', label='Polynomial Regression', linestyle='--', color='blue')
plt.show()

# Fit a seasonal curve to the seasonal component
period_coef_ = seasonal_curve_fit_and_plot(f=seasonal_curve_function, x=df['timestamp_sec'], y=result.seasonal)

plt.figure(figsize=(20, 20))
plot_time_series(df['timestamp_sec'], result.seasonal, label='Base', title='Base', linestyle='-', color='red')
plot_time_series(df['timestamp_sec'],
                 seasonal_curve_function(df['timestamp_sec'], period_coef_[0], period_coef_[1], period_coef_[2],
                                         period_coef_[3], ), label='Fitted', title='Curve Fit', linestyle='--',
                 color='blue')
plt.show()

# Complete the fitting process by combining polynomial and seasonal fits
y_fitted_, fitting_error_, fitted_function_ = complete_fit_and_plot(df["timestamp_sec"], df["value"], poly_coef_,
                                                                    period_coef_)

# Plot the base trend, fitted trend, and fitting error
plt.figure(figsize=(20, 20))
# plt.subplot(3, 1, 1)
plot_time_series(df["timestamp_sec"], df["value"], title='Base', label='Base', linestyle='-', color='red')
# plt.subplot(3, 1, 2)
plot_time_series(df["timestamp_sec"], y_fitted_, title='Fitted', label='Fitted', linestyle='--', color='blue')
# plt.subplot(3, 1, 3)
plot_time_series(df["timestamp_sec"], fitting_error_, title='Error', label='Error', linestyle='-', color='green')
plt.show()

# Check the quality of fitting and print metrics
check_fitting_quality_and_print_metrics(df["value"], y_fitted_)
e_mean, e_std = print_error_distribution_and_return_stats(fitting_error_)

# Visualize the distribution of errors
plt.figure(figsize=(20, 20))
sns.histplot(fitting_error_, bins=30, kde=True, color='green')
plt.title('Distribution of Errors')
plt.xlabel('Error')
plt.ylabel('Frequency')
plt.show()

serialized_function = dill.dumps(fitted_function_)
deserialized_function = dill.loads(serialized_function)

# Generate x values
x_values = np.linspace(x_lower_limit, x_upper_limit, 1000, dtype=int)
print(x_values)
# Calculate the trend values
trend_values = deserialized_function(x_values)
print(trend_values)
print(x_lower_limit, "mean:", deserialized_function(np.linspace(x_lower_limit, x_upper_limit, 2)))
print("mean:", deserialized_function(x_upper_limit))
# Calculate the maximum and minimum trend values
max_trend_value = np.max(trend_values)
min_trend_value = np.min(trend_values)
# Set the lower and upper limits of y based on trend values and standard deviation
y_lower_limit = min_trend_value - 4 * e_std
y_upper_limit = max_trend_value + 4 * e_std
print("y_lower_limit:", y_lower_limit, "y_upper_limit:", y_upper_limit)
# Generate the surface function using the trend and standard deviation
surface = generate_surface_function(e_std, deserialized_function)
print("mean:", deserialized_function(x_lower_limit), "prob:",
      surface(deserialized_function(x_lower_limit), x_lower_limit))
print("mean:", deserialized_function(x_upper_limit), "prob:",
      surface(deserialized_function(x_upper_limit), x_upper_limit))
# Calculate the probability of y being out of bounds
prob = calculate_probability(x_lower_limit, x_upper_limit, y_lower_limit, y_upper_limit, surface, y_lower_bound)

# Create a 3D plot figure with two subplots
fig = plt.figure(figsize=(20, 20))

# Plot the surface for the specified range
plot_surface(fig, x_lower_limit, x_upper_limit, y_lower_limit, y_upper_limit, surface, 211, y_lower_limit,
             y_upper_limit, )

# Plot another surface for a different range in the second subplot
plot_surface(fig, x_lower_limit, x_upper_limit, y_lower_bound, y_upper_limit, surface, 212, y_lower_limit,
             y_upper_limit, )

# Show the plot
plt.show()
