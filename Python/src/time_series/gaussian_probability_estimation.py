import numpy as np
from scipy.stats import norm


def calculate_probability(mean, std_dev, threshold):
    """
    Calculate the probability that a random variable exceeds a given threshold.

    Parameters:
    - mean: Mean of the random variable.
    - std_dev: Standard deviation of the random variable.
    - threshold: Threshold value.

    Returns:
    - The probability that the random variable is greater than the threshold, in percentage.
    """
    # Calculate the z-score
    z_score = (threshold - mean) / std_dev
    # Calculate the probability using the complementary cumulative distribution function (CDF)
    probability = 1 - norm.cdf(z_score)
    # Convert the probability to a percentage and round it
    return probability * 100


def evaluate_points(function, lower_bound, upper_bound, num_points=100):
    """
    Generate the mean values within the specified range using the given function and number of points.

    Parameters:
    - function: The function to evaluate.
    - lower_bound: Lower bound of the range.
    - upper_bound: Upper bound of the range.
    - num_points: Number of points to generate within the range.

    Returns:
    - mean_values: List of mean values evaluated from the function.
    - base_values: List of base values (independent variable) within the range.
    """
    base_values = np.linspace(lower_bound, upper_bound, num_points)
    return [function(x) for x in base_values], base_values


def calculate_mean_probability(function, std_dev, lower_bound, upper_bound, threshold, num_points=100):
    """
    Calculate the mean probability that a function of a random variable exceeds a threshold within a given range.

    Parameters:
    - function: The function representing the relationship between the random variable and another variable.
    - std_dev: Standard deviation of the random variable.
    - lower_bound: Lower bound of the range for the independent variable.
    - upper_bound: Upper bound of the range for the independent variable.
    - threshold: Threshold value.
    - num_points: Number of points used for approximation within the range.

    Returns:
    - mean_probability: The mean probability that the function of the random variable exceeds the threshold within
    the given range.
    - probabilities: List of probabilities corresponding to each mean value within the range.
    - mean_values: List of mean values within the range [lower_bound, upper_bound].
    - base_values: List of base values (independent variable) within the range.
    """
    # Generate the mean values within the range [lower_bound, upper_bound]
    mean_values, base_values = evaluate_points(function, lower_bound, upper_bound, num_points)
    # Calculate probabilities for each x in the range [lower_bound, upper_bound]
    probabilities = [calculate_probability(x, std_dev, threshold) for x in mean_values]
    # Calculate the mean probability
    mean_probability = round(sum(probabilities) / len(probabilities), 2)
    print(f"Mean probability of y > {threshold} given {lower_bound} < x < {upper_bound}: {mean_probability}%")
    return mean_probability, probabilities, mean_values, base_values
