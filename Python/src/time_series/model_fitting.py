import numpy as np
import pandas as pd
from scipy.optimize import curve_fit
from scipy.stats import norm
from sklearn.linear_model import LinearRegression
from sklearn.metrics import mean_squared_error, r2_score
from sklearn.preprocessing import PolynomialFeatures
from statsmodels.tsa.filters.hp_filter import hpfilter
from statsmodels.tsa.seasonal import seasonal_decompose


def format_data(data):
    # Convert the list of tuples into a pandas DataFrame
    df = pd.DataFrame(data, columns=['timestamp_sec', 'value'])
    # Convert the 'value' column to float
    df['value'] = df['value'].astype(float)
    return df


# Function to apply the Hodrick-Prescott filter with optimal lambda
def apply_hp_filter_with_optimal_lambda(time_series, lambda_range):
    # Initialize variables to store the best results
    best_lambda = None
    best_trend = None
    best_cycle = None
    best_variance_ratio = float('inf')  # Initialize with a large value

    # Loop through different lambda values
    for lambda_ in lambda_range:
        # Apply Hodrick-Prescott filter
        cycle_, trend_ = hpfilter(time_series, lamb=lambda_)

        # Calculate variance ratio
        variance_ratio = np.var(cycle_) / np.var(time_series)

        # Check if the current result is better than the previous best
        if variance_ratio < best_variance_ratio:
            best_lambda = lambda_
            best_trend = trend_
            best_cycle = cycle_
            best_variance_ratio = variance_ratio

    return best_cycle, best_trend, best_lambda


# Function for seasonal decomposition
def seasonal_decomposition(time_series, period_range):
    # Initialize variables to store the best results
    best_period = None
    best_result = None
    best_seasonal_variance = float('-inf')  # Initialize with a small value
    best_residual_variance = float('inf')  # Initialize with a large value

    # Loop through different period values
    for period in period_range:
        # Apply seasonal decomposition
        result_ = seasonal_decompose(time_series, model='additive', period=period)

        # Calculate seasonal and residual variances
        seasonal_variance = np.var(result_.seasonal)
        residual_variance = np.var(result_.resid)

        # Check if the current result is better than the previous best
        if seasonal_variance > best_seasonal_variance and residual_variance < best_residual_variance:
            best_period = period
            best_result = result_
            best_seasonal_variance = seasonal_variance
            best_residual_variance = residual_variance

    print(f"Optimal period: {best_period}")
    return best_result


# Function for polynomial fitting
def polynomial_fit_and_plot(x, y, degree=1):
    pr = PolynomialFeatures(degree=degree)
    x_poly = pr.fit_transform(x)
    lr_2 = LinearRegression()
    lr_2.fit(x_poly, y)
    # Get the coefficients
    coefficients = lr_2.coef_[0]
    print("Polynomial coefficients :", coefficients)
    return coefficients


# Function for curve fitting using scipy.optimize.curve_fit
def seasonal_curve_fit_function(x, a, b, c):
    return a * x + b * np.sin(c * x)


def seasonal_curve_fit_and_plot(f, x, y):
    res = curve_fit(f, x, y)
    # Extract the fitted parameters
    a, b, c = res[0]
    print("Curve fit coefficients :", res[0])
    return a, b, c


# Function for complete fitting
def complete_fit_and_plot(x, y, poly_coef, period_coef):
    def fitted_function(x_): return poly_coef[0] + poly_coef[1] * x_ + poly_coef[2] * x_ ** 2 + period_coef[0] * x_ + \
        period_coef[1] * np.sin(period_coef[2] * x_)

    y_fitted = np.squeeze(fitted_function(x))
    fitting_error = y - y_fitted
    return y_fitted, fitting_error, fitted_function


# Function to check fitting quality
def check_fitting_quality_and_print_metrics(y, y_fitted):
    mse = mean_squared_error(y, y_fitted)
    rmse = np.sqrt(mse)
    r_squared = r2_score(y, y_fitted)
    print(f'Mean Squared Error: {mse}')
    print(f'Root Mean Squared Error: {rmse}')
    print(f'R-squared: {r_squared}')


# Function to print error distribution
def print_error_distribution_and_return_stats(fitting_error):
    error_mean, error_std = norm.fit(fitting_error)
    print(f"Mean of the error: {error_mean}")
    print(f"Standard deviation of the error: {error_std}")
    return error_mean, error_std


def reevaluate_model(data):
    df = format_data(data)
    # Apply Hodrick-Prescott filter with optimal lambda
    cycle, trend, optimal_lambda = apply_hp_filter_with_optimal_lambda(time_series=df['value'],
                                                                       lambda_range=[1, 10, 50, 100, 200, 500, 0.5])
    # Perform seasonal decomposition on the trend component
    result = seasonal_decomposition(time_series=trend, period_range=[10, 20, 50, 100, 200, 500, 1000, 2000])
    df_trend = df.copy()
    df_trend["value"] = result.trend
    df_trend = df_trend.dropna()
    # Use polynomial_fit_and_plot to fit a polynomial to the trend component
    poly_coef_ = polynomial_fit_and_plot(x=df_trend["timestamp_sec"].values.reshape(-1, 1),
                                         y=df_trend["value"].values.reshape(-1, 1), degree=2)
    # Use seasonal_curve_fit_and_plot to fit the function to the seasonal component
    period_coef_ = seasonal_curve_fit_and_plot(f=seasonal_curve_fit_function, x=df['X'], y=result.seasonal)
    # Compose the fitted trend and seasonal component
    y_fitted_, fitting_error_, fitted_function_ = complete_fit_and_plot(df["timestamp_sec"], df["value"], poly_coef_,
                                                                        period_coef_)
    # Check the quality of the fitting using the fitting error
    check_fitting_quality_and_print_metrics(df["value"], y_fitted_)
    # Fit the error distribution with a gaussian
    e_mean, e_std = print_error_distribution_and_return_stats(fitting_error_)
    return fitted_function_, e_std
