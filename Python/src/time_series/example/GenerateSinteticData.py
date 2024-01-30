import numpy as np
import pandas as pd

# Generate synthetic data
X = np.linspace(start=0, stop=200, num=4000).reshape(-1, 1)
true_function = np.squeeze(0.01 * X + np.sin(X))
noise = np.random.normal(0, 0.1, size=X.shape[0])
Y = true_function + noise

# Create a DataFrame
df = pd.DataFrame({'X': np.squeeze(X), 'Y': Y})

# Save DataFrame to a CSV file
df.to_csv('../../data/generated_data.csv', index=False)
