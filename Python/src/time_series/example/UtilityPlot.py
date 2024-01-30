import numpy as np
from matplotlib import pyplot as plt


def plot_time_series(x, y, label=None, xlabel='X', ylabel='Values', title=None, linestyle='-', color='blue'):
    plt.plot(x, y, label=label, linestyle=linestyle, color=color)
    plt.xlabel(xlabel)
    plt.ylabel(ylabel)
    if title:
        plt.title(title)
    plt.legend()


def plot_surface(fig, x_l, x_u, y_l, y_u, surf, position, y_b_l=-20, y_b_u=20, z_u=1):
    """
    Plot a 3D surface plot.

    Parameters:
    - fig (Figure): Figure object.
    - x_l (float): Lower limit of x.
    - x_u (float): Upper limit of x.
    - y_l (float): Lower limit of y.
    - y_u (float): Upper limit of y.
    - surf (function): Surface function.
    - position (int): Position of subplot.
    - y_b_l (float): Lower bound of y.
    - y_b_u (float): Upper bound of y.
    - z_u (float): Upper bound of z.

    Returns:
    - None
    """
    ax = fig.add_subplot(position, projection='3d')

    # Create a set of data points for x and y
    x = np.linspace(x_l, x_u, 1000)
    y = np.linspace(y_l, y_u, 1000)
    x, y = np.meshgrid(x, y)

    # Compute the values of the surface function
    z = surf(y=y, x=x)

    # Plot the 3D surface with customizations
    ax.plot_surface(x, y, z, cmap='viridis')
    ax.view_init(elev=60, azim=0)  # Set the viewing angle
    ax.set_ylim([y_b_l, y_b_u])
    ax.set_zlim([0, z_u])

    # Add labels to the axes
    ax.set_xlabel('X')
    ax.set_ylabel('Y')
    ax.set_zlabel('Z')
