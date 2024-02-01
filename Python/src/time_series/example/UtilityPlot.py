import numpy as np
from matplotlib import pyplot as plt
from scipy.stats import norm


def plot_time_series(x, y, label=None, xlabel='X', ylabel='Values', title=None, linestyle='-', color='blue'):
    plt.plot(x, y, label=label, linestyle=linestyle, color=color)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    if title:
        plt.title(title)
    plt.legend()


def display_probability_surface(base_values, function, std_dev, threshold):
    """
    Display the probability surface and function curve.

    Args:
    - base_values: array-like, the base values of the independent variable
    - function: function, the function to visualize
    - std_dev: float, standard deviation for the normal distribution
    - threshold: float, threshold value for the probability distribution
    """

    # Independent variable
    x = base_values

    # Determine plot ranges
    y_lower_limit, y_upper_limit = min(min(function(x)) - 2 * std_dev, threshold - std_dev), max(
        max(function(x)) + 2 * std_dev, threshold + std_dev)

    # Values for dependent variable
    y = np.linspace(y_lower_limit, y_upper_limit, 100)

    # Create the figure
    fig = plt.figure(figsize=(20, 10))

    # Subplot 1 for the function curve
    ax1 = fig.add_subplot(121)
    ax1.plot(x, function(x), label='Function Curve')
    ax1.axhline(y=threshold, color='r', linestyle='--', label='Threshold')
    ax1.set_xlabel('Independent Variable (x)')
    ax1.set_ylabel('Dependent Variable (y)')
    ax1.set_title('Function Curve')
    ax1.set_ylim([y_lower_limit, y_upper_limit])
    ax1.grid(True)
    ax1.legend()

    # Create a meshgrid for 3D plotting
    X, Y = np.meshgrid(x, y)

    # Compute the normal probability density function
    Z1 = norm.pdf(Y, loc=function(X), scale=std_dev)

    # Create a copy for areas where y <= threshold
    Z2 = np.copy(Z1)
    Z2[Y <= threshold] = np.nan

    # Subplot 2 for the comparison of surfaces
    ax2 = fig.add_subplot(122, projection='3d')
    ax2.view_init(elev=80, azim=-90)
    ax2.plot_surface(X, Y, Z1, cmap='viridis', alpha=0.7, label='Original Surface')
    ax2.plot_surface(X, Y, Z2, color='red', alpha=1, label='Y > Threshold Surface')
    ax2.set_xlabel('Independent Variable (x)')
    ax2.set_ylabel('Dependent Variable (y)')
    ax2.set_zlabel('PDF')
    ax2.set_title('Comparison of Surfaces')
    ax2.set_ylim([y_lower_limit, y_upper_limit])
    ax2.legend()

    # Show the plot
    plt.show()
