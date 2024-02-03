import numpy as np
from scipy.stats import norm


def calculate_probability(mean, std_dev, threshold_l=float('-inf'), threshold_u=float('inf')):
    """
    Calculate the probability that a random variable falls outside a given range.

    Parameters:
    - mean: Mean of the random variable.
    - std_dev: Standard deviation of the random variable.
    - threshold_l: Lower threshold value.
    - threshold_u: Upper threshold value.

    Returns:
    - The probability that the random variable falls outside the range defined by threshold_l and threshold_u,
    in percentage.
    """
    # Calculate the z-scores
    z_score_l = (threshold_l - mean) / std_dev
    z_score_u = (threshold_u - mean) / std_dev

    # Calculate the probability using the complementary cumulative distribution function (CDF)
    probability_within_range = norm.cdf(z_score_u) - norm.cdf(z_score_l)

    # Calculate the probability of falling outside the range
    probability_outside_range = 1 - probability_within_range

    # Convert the probability to a percentage and round it
    return probability_outside_range * 100


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


def calculate_mean_probability(function, std_dev, lower_bound, upper_bound, threshold_l=float('-inf'),
                               threshold_u=float('inf'), num_points=100):
    """
    Calculate the mean probability that a function of a random variable exceeds a threshold within a given range.

    Parameters:
    - function: The function representing the relationship between the random variable and another variable.
    - std_dev: Standard deviation of the random variable.
    - lower_bound: Lower bound of the range for the independent variable.
    - upper_bound: Upper bound of the range for the independent variable.
    - threshold_l: Lower threshold value.
    - threshold_u: Upper threshold value.
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
    probabilities = [calculate_probability(x, std_dev, threshold_u=threshold_u, threshold_l=threshold_l) for x in
                     mean_values]
    # Calculate the mean probability
    mean_probability = round(sum(probabilities) / len(probabilities), 2)
    print(f"Mean probability of y > {threshold_u} or y < {threshold_l} given {lower_bound} < x < {upper_bound}: "
          f"{mean_probability}%")
    return mean_probability, probabilities, mean_values, base_values
