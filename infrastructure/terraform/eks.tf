resource "aws_iam_role" "eks_iam" {
  name = "example_go_api"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "eks.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

resource "aws_iam_role_policy_attachment" "eks_iam-AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_iam.name
}

resource "aws_eks_cluster" "eks_cluster" {
  name     = "example_go_api"
  role_arn = aws_iam_role.eks_iam.arn

  vpc_config {
    subnet_ids = [
      aws_subnet.private_1a.id,
      aws_subnet.private_1b.id,
      aws_subnet.public_1a.id,
      aws_subnet.public_1b.id
    ]
  }

  depends_on = [aws_iam_role_policy_attachment.eks_iam-AmazonEKSClusterPolicy]
}
